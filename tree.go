/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-15 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 13:55:55
 * @FilePath: \gosh\tree.go
 * @Description: 路由树的实现，用于处理 URL 路径和参数。
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/gosh/constants"
)

// Param 表示单个 URL 参数，由键和值组成
type Param struct {
	Key   string // 参数名
	Value string // 参数值
}

// Params 是一个 Param 切片，由路由器返回
// 切片是有序的，第一个 URL 参数也是第一个切片值
type Params []Param

// Get 获取指定名称的路径参数
func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

// ByName 根据参数名称返回第一个匹配的参数值
// 如果没有找到匹配的参数，则返回空字符串
func (ps Params) ByName(name string) string {
	value, _ := ps.Get(name)
	return value
}

// methodTree 结构体表示 HTTP 方法与其对应的路由树
type methodTree struct {
	method string // HTTP 方法
	root   *Node  // 路由树的根节点
}

// methodTrees 是 methodTree 的切片
type methodTrees []methodTree

// get 根据 HTTP 方法获取对应的路由树的根节点
func (trees methodTrees) get(method string) *Node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

// addChild 将子节点添加到当前节点，保持通配符子节点在最后
func (n *Node) addChild(child *Node) {
	if n.wildChild && len(n.children) > 0 {
		wildcardChild := n.children[len(n.children)-1]
		n.children = append(n.children[:len(n.children)-1], child, wildcardChild)
	} else {
		n.children = append(n.children, child)
	}
}

// nodeType 表示节点的类型
type nodeType uint8

const (
	staticNode   nodeType = iota // 静态节点
	rootNode                     // 根节点
	paramNode                    // 参数节点
	wildcardNode                 // 通配符节点
)

// Node 表示树中的节点
type Node struct {
	path      string        // 节点的路径
	indices   string        // 子节点的索引字符
	wildChild bool          // 是否有通配符子节点
	nType     nodeType      // 节点类型
	priority  uint32        // 节点的优先级
	children  []*Node       // 子节点，最多有一个参数节点在数组的末尾
	handlers  HandlersChain // 处理函数链
	fullPath  string        // 完整路径
}

// incrementChildPrio 增加给定子节点的优先级，并在必要时重新排序
func (n *Node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// 调整位置（移动到前面）
	newPos := pos
	for newPos > 0 && cs[newPos-1].priority < prio {
		// 交换节点位置
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	// 构建新的索引字符字符串
	if newPos != pos {
		n.indices = n.indices[:newPos] + // 不变前缀，可能为空
			n.indices[pos:pos+1] + // 移动的索引字符
			n.indices[newPos:pos] + n.indices[pos+1:] // 剩余部分
	}

	return newPos
}

// addRoute 添加一个节点到路径中，并注册处理函数
// 不是线程安全的！
func (n *Node) addRoute(path string, handlers HandlersChain) {
	fullPath := path
	n.priority++

	// 如果是空树
	if len(n.path) == 0 && len(n.children) == 0 {
		n.insertChild(path, fullPath, handlers)
		n.nType = rootNode
		return
	}

	parentFullPathIndex := 0

walk:
	for {
		// 找出最长的公共前缀
		i := mathx.LongestCommonPrefix(path, n.path)

		// 如果找到公共前缀，进行节点分割
		if i < len(n.path) {
			n.splitNode(i, path, fullPath)
		}

		// 处理路径的剩余部分
		if i < len(path) {
			path = path[i:]
			c := path[0]

			// 处理参数节点的特殊情况
			if n.nType == paramNode && c == constants.PathSeparator && len(n.children) == 1 {
				parentFullPathIndex += len(n.path)
				n = n.children[0]
				n.priority++
				continue walk
			}

			// 处理子节点
			if n.handleChildNode(c, path, fullPath, handlers, &parentFullPathIndex) {
				return
			}
		}

		// 注册处理函数
		n.registerHandlers(fullPath, handlers)
		return
	}
}

// splitNode 分割节点，将当前节点的路径与子节点分开
func (n *Node) splitNode(i int, path string, fullPath string) {
	child := Node{
		path:      n.path[i:], // 子节点的路径
		wildChild: n.wildChild,
		nType:     staticNode,
		indices:   n.indices,
		children:  n.children,
		handlers:  n.handlers,
		priority:  n.priority - 1,
		fullPath:  n.fullPath,
	}

	n.children = []*Node{&child}                             // 设置子节点
	n.indices = convert.SliceByteToString([]byte{n.path[i]}) // 更新索引
	n.path = path[:i]                                        // 更新当前节点的路径
	n.handlers = nil                                         // 清空当前节点的处理函数
	n.wildChild = false
	n.fullPath = fullPath[:len(fullPath)-len(path)+i] // 更新全路径
}

// handleChildNode 处理子节点的逻辑
func (n *Node) handleChildNode(c byte, path string, fullPath string, handlers HandlersChain, parentFullPathIndex *int) bool {
	// 查找具有相同路径字节的子节点
	for i, maxIndices := 0, len(n.indices); i < maxIndices; i++ {
		if c == n.indices[i] {
			*parentFullPathIndex += len(n.path)
			return false
		}
	}

	// 如果当前字符不是路径参数前缀或通配符
	if c != constants.PathParamPrefix && c != constants.WildcardSymbol && n.nType != wildcardNode {
		n.indices += convert.SliceByteToString([]byte{c}) // 更新索引
		child := &Node{
			fullPath: fullPath, // 设置子节点的全路径
		}
		n.addChild(child)                        // 添加子节点
		n.incrementChildPrio(len(n.indices) - 1) // 增加子节点优先级
		n = child                                // 进入新创建的子节点
	} else if n.wildChild {
		n = n.children[len(n.children)-1] // 进入最后一个通配符子节点
		n.priority++                      // 增加优先级

		// 检查路径是否与当前通配符匹配
		if len(path) >= len(n.path) && n.path == path[:len(n.path)] &&
			n.nType != wildcardNode &&
			(len(n.path) >= len(path) || path[len(n.path)] == constants.PathSeparator) {
			return false
		}

		// 检查通配符冲突
		n.checkWildcardConflict(path, fullPath)
	}

	n.insertChild(path, fullPath, handlers) // 插入新的子节点
	return true
}

// checkWildcardConflict 检查通配符冲突
func (n *Node) checkWildcardConflict(path string, fullPath string) {
	pathSeg := path
	if n.nType != wildcardNode {
		pathSeg = strings.SplitN(pathSeg, constants.PathSeparatorStr, 2)[0] // 取路径段
	}
	prefix := fullPath[:strings.Index(fullPath, pathSeg)] + n.path // 构建前缀
	panic("'" + pathSeg +
		"' in new path '" + fullPath +
		"' conflicts with existing wildcard '" + n.path +
		"' in existing prefix '" + prefix +
		"'")
}

// registerHandlers 注册处理函数
func (n *Node) registerHandlers(fullPath string, handlers HandlersChain) {
	if n.handlers != nil {
		panic("handlers are already registered for path '" + fullPath + "'") // 检查是否已经注册处理函数
	}
	n.handlers = handlers // 注册处理函数
	n.fullPath = fullPath // 设置全路径
}

// findWildcard 查找路径中的通配符段，并检查名称是否包含无效字符
// 如果没有找到通配符，则返回 -1 作为索引
func findWildcard(path string) (wildcard string, i int, valid bool) {
	// 查找起始
	for start, c := range []byte(path) {
		if c != constants.PathParamPrefix && c != constants.WildcardSymbol {
			continue
		}

		// 查找结尾并检查无效字符
		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case constants.PathSeparator:
				return path[start : start+1+end], start, valid
			case constants.PathParamPrefix, constants.WildcardSymbol:
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}

// insertChild 插入子节点
func (n *Node) insertChild(path string, fullPath string, handlers HandlersChain) {
	for {
		// 查找通配符
		wildcard, i, valid := findWildcard(path)
		if i < 0 { // 没有找到通配符
			break
		}

		// 通配符名称只能包含一个 ':' 或 '*' 字符
		validateWildcard(wildcard, fullPath, valid)

		if wildcard[0] == constants.PathParamPrefix { // paramNode
			if i > 0 {
				// 在当前通配符之前插入前缀
				n.path = path[:i]
				path = path[i:]
			}

			child := &Node{
				nType:    paramNode,
				path:     wildcard,
				fullPath: fullPath,
			}
			n.addChild(child)
			n.wildChild = true
			n = child
			n.priority++

			// 如果路径不以通配符结尾，则会有另一个子路径以 constants.Separator 开始
			if len(wildcard) < len(path) {
				path = path[len(wildcard):]

				child := &Node{
					priority: 1,
					fullPath: fullPath,
				}
				n.addChild(child)
				n = child
				continue
			}

			// 否则我们完成了。在新叶中插入处理函数
			n.handlers = handlers
			return
		}

		// wildcardNode
		if i+len(wildcard) != len(path) {
			panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
		}

		if len(n.path) > 0 && n.path[len(n.path)-1] == constants.PathSeparator {
			pathSeg := strings.SplitN(n.children[0].path, constants.PathSeparatorStr, 2)[0]
			panic("catch-all wildcard '" + path +
				"' in new path '" + fullPath +
				"' conflicts with existing path segment '" + pathSeg +
				"' in existing prefix '" + n.path + pathSeg +
				"'")
		}

		// 当前固定宽度为 1
		i--
		if path[i] != constants.PathSeparator {
			panic("no / before catch-all in path '" + fullPath + "'")
		}

		n.path = path[:i]

		// 第一个节点：具有空路径的 wildcardNode 节点
		child := &Node{
			wildChild: true,
			nType:     wildcardNode,
			fullPath:  fullPath,
		}

		n.addChild(child)
		n.indices = string(constants.PathSeparator)
		n = child
		n.priority++

		// 第二个节点：持有变量的节点
		child = &Node{
			path:     path[i:],
			nType:    wildcardNode,
			handlers: handlers,
			priority: 1,
			fullPath: fullPath,
		}
		n.children = []*Node{child}

		return
	}

	// 如果没有找到通配符，简单地插入路径和处理函数
	n.path = path
	n.handlers = handlers
	n.fullPath = fullPath
}

// validateWildcard 验证通配符的有效性
func validateWildcard(wildcard string, fullPath string, valid bool) {
	// 通配符名称只能包含一个 ':' 或 '*' 字符
	if !valid {
		panic("only one wildcard per path segment is allowed, has: '" +
			wildcard + "' in path '" + fullPath + "'")
	}
	if len(wildcard) < 2 {
		panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
	}
}

// nodeValue 保存 (*Node).getValue 方法的返回值
type nodeValue struct {
	handlers HandlersChain // 处理函数链
	params   *Params       // 路径参数
	tsr      bool          // 是否为尾随斜杠
	fullPath string        // 完整路径
}

// skippedNode 用于存储跳过的节点信息
type skippedNode struct {
	path        string // 路径
	node        *Node  // 节点
	paramsCount int16  // 参数计数
}

// getValue 根据给定路径（键）返回已注册的处理函数
// 通配符的值保存到一个映射中
// 如果找不到处理函数，则根据给定路径的额外（无）尾随斜杠推荐 TSR（尾随斜杠重定向）
func (n *Node) getValue(path string, params *Params, skippedNodes *[]skippedNode) (value nodeValue) {
	var globalParamsCount int16

walk: // 外部循环遍历树
	for {
		prefix := n.path
		if len(path) > len(prefix) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]

				// 首先尝试匹配所有非通配符子节点
				idxc := path[0]
				for i, c := range []byte(n.indices) {
					if c == idxc {
						// 处理通配符子节点，始终在数组的最后
						if n.wildChild {
							index := len(*skippedNodes)
							*skippedNodes = (*skippedNodes)[:index+1]
							(*skippedNodes)[index] = skippedNode{
								path: prefix + path,
								node: &Node{
									path:      n.path,
									wildChild: n.wildChild,
									nType:     n.nType,
									priority:  n.priority,
									children:  n.children,
									handlers:  n.handlers,
									fullPath:  n.fullPath,
								},
								paramsCount: globalParamsCount,
							}
						}

						n = n.children[i]
						continue walk
					}
				}

				if !n.wildChild {
					// 如果路径在循环结束时不等于 '/' 并且当前节点没有子节点
					// 当前节点需要回滚到最后一个有效的 skippedNode
					if path != constants.PathSeparatorStr {
						for length := len(*skippedNodes); length > 0; length-- {
							skippedNode := (*skippedNodes)[length-1]
							*skippedNodes = (*skippedNodes)[:length-1]
							if strings.HasSuffix(skippedNode.path, path) {
								path = skippedNode.path
								n = skippedNode.node
								if value.params != nil {
									*value.params = (*value.params)[:skippedNode.paramsCount]
								}
								globalParamsCount = skippedNode.paramsCount
								continue walk
							}
						}
					}

					// 没有找到
					// 如果路径的末尾是 constants.PathSeparatorStr 且当前节点是叶子节点
					// 则可以推荐重定向到同一 URL
					value.tsr = path == constants.PathSeparatorStr && n.handlers != nil
					return
				}

				// 处理通配符子节点，始终在数组的最后
				n = n.children[len(n.children)-1]
				globalParamsCount++

				switch n.nType { //nolint:exhaustive
				case paramNode:
					// 修复参数的截断
					// tree_test.go  line: 204

					// 查找 paramNode 结束（要么是 constants.PathSeparator 要么是路径结束）
					end := 0
					for end < len(path) && path[end] != constants.PathSeparator {
						end++
					}

					// 保存 paramNode 值
					if params != nil && cap(*params) > 0 {
						if value.params == nil {
							value.params = params
						}
						// 在预分配的容量内扩展切片
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						val := path[:end]
						(*value.params)[i] = Param{
							Key:   n.path[1:], // 去掉前导冒号
							Value: val,
						}
					}

					// 我们需要更深入！
					if end < len(path) {
						if len(n.children) > 0 {
							path = path[end:]
							n = n.children[0]
							continue walk
						}

						// ...但我们不能
						value.tsr = len(path) == end+1
						return
					}

					if value.handlers = n.handlers; value.handlers != nil {
						value.fullPath = n.fullPath
						return
					}
					if len(n.children) == 1 {
						// 没有找到处理函数。检查路径 + 尾随斜杠是否存在处理函数
						n = n.children[0]
						value.tsr = (n.path == constants.PathSeparatorStr && n.handlers != nil) || (n.path == "" && n.indices == constants.PathSeparatorStr)
					}
					return

				case wildcardNode:
					// 保存 paramNode 值
					if params != nil {
						if value.params == nil {
							value.params = params
						}
						// 在预分配的容量内扩展切片
						i := len(*value.params)
						*value.params = (*value.params)[:i+1]
						(*value.params)[i] = Param{
							Key:   n.path[1:], // 去掉前导通配符
							Value: path,       // 通配符的值
						}
					}
					value.handlers = n.handlers
					value.fullPath = n.fullPath
					return

				default:
					panic("无效的节点类型")
				}
			}
		}

		if path == prefix {
			// 如果当前路径不等于 constants.PathSeparatorStr 并且节点没有注册的处理函数
			// 且最近匹配的节点有子节点，则需要回滚到最后一个有效的 skippedNode
			if n.handlers == nil && path != constants.PathSeparatorStr {
				for length := len(*skippedNodes); length > 0; length-- {
					skippedNode := (*skippedNodes)[length-1]
					*skippedNodes = (*skippedNodes)[:length-1]
					if strings.HasSuffix(skippedNode.path, path) {
						path = skippedNode.path
						n = skippedNode.node
						if value.params != nil {
							*value.params = (*value.params)[:skippedNode.paramsCount]
						}
						globalParamsCount = skippedNode.paramsCount
						continue walk
					}
				}
			}

			// 我们应该到达包含处理函数的节点
			// 检查此节点是否注册了处理函数
			if value.handlers = n.handlers; value.handlers != nil {
				value.fullPath = n.fullPath
				return
			}

			// 如果此路径没有处理函数，但此路径有通配符子节点
			// 则必须存在处理此路径的附加尾随斜杠的处理函数
			if path == constants.PathSeparatorStr && n.wildChild && n.nType != rootNode {
				value.tsr = true
				return
			}

			if path == constants.PathSeparatorStr && n.nType == staticNode {
				value.tsr = true
				return
			}

			// 没有找到处理函数。检查路径 + 尾随斜杠是否存在处理函数
			for i, c := range []byte(n.indices) {
				if c == constants.PathSeparator {
					n = n.children[i]
					value.tsr = (len(n.path) == 1 && n.handlers != nil) ||
						(n.nType == wildcardNode && n.children[0].handlers != nil)
					return
				}
			}
			return
		}

		// 没有找到。我们可以推荐重定向到同一 URL，
		// 通过增加额外的尾随斜杠
		value.tsr = path == constants.PathSeparatorStr ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == constants.PathSeparator &&
				path == prefix[:len(prefix)-1] && n.handlers != nil)

		// 回滚到最后一个有效的 skippedNode
		if !value.tsr && path != constants.PathSeparatorStr {
			for length := len(*skippedNodes); length > 0; length-- {
				skippedNode := (*skippedNodes)[length-1]
				*skippedNodes = (*skippedNodes)[:length-1]
				if strings.HasSuffix(skippedNode.path, path) {
					path = skippedNode.path
					n = skippedNode.node
					if value.params != nil {
						*value.params = (*value.params)[:skippedNode.paramsCount]
					}
					globalParamsCount = skippedNode.paramsCount
					continue walk
				}
			}
		}
		return
	}
}

// shiftNRuneBytes 将数组中的字节左移 n 个字节
func shiftNRuneBytes(rb [4]byte, n int) [4]byte { //nolint:unused
	switch n {
	case 0:
		return rb
	case 1:
		return [4]byte{rb[1], rb[2], rb[3], 0}
	case 2:
		return [4]byte{rb[2], rb[3]}
	case 3:
		return [4]byte{rb[3]}
	default:
		return [4]byte{}
	}
}

// findCaseInsensitivePathRec 递归查找路径，用于不区分大小写的路径查找
func (n *Node) findCaseInsensitivePathRec(path string, ciPath []byte, rb [4]byte, fixTrailingSlash bool) []byte { //nolint:unused
	npLen := len(n.path)

walk: // 外部循环遍历树
	for len(path) >= npLen && (npLen == 0 || strings.EqualFold(path[1:npLen], n.path[1:])) {
		// 将公共前缀添加到结果
		oldPath := path
		path = path[npLen:]
		ciPath = append(ciPath, n.path...)

		if len(path) == 0 {
			// 我们应该到达包含处理函数的节点
			// 检查此节点是否注册了处理函数
			if n.handlers != nil {
				return ciPath
			}

			// 没有找到处理函数
			// 尝试通过添加尾随斜杠来修复路径
			if fixTrailingSlash {
				for i, c := range []byte(n.indices) {
					if c == constants.PathSeparator {
						n = n.children[i]
						if (len(n.path) == 1 && n.handlers != nil) ||
							(n.nType == wildcardNode && n.children[0].handlers != nil) {
							return append(ciPath, constants.PathSeparator)
						}
						return nil
					}
				}
			}
			return nil
		}

		// 如果此节点没有通配符（参数节点或通配符节点）子节点
		// 我们可以直接查找下一个子节点并继续遍历树
		if !n.wildChild {
			// 跳过已经处理的 rune 字节
			rb = shiftNRuneBytes(rb, npLen)

			if rb[0] != 0 {
				// 旧的 rune 尚未处理完
				idxC := rb[0]
				for i, c := range []byte(n.indices) {
					if c == idxC {
						// 继续处理子节点
						n = n.children[i]
						npLen = len(n.path)
						continue walk
					}
				}
			} else {
				// 处理新的 rune
				var rv rune

				// 查找 rune 开始位置
				// Runes 最多 4 字节长
				var off int
				for maxNPLen := mathx.AtMost(npLen, 3); off < maxNPLen; off++ {
					if i := npLen - off; utf8.RuneStart(oldPath[i]) {
						// 从缓存路径中读取 rune
						rv, _ = utf8.DecodeRuneInString(oldPath[i:])
						break
					}
				}

				// 计算当前 rune 的小写字节
				lo := unicode.ToLower(rv)
				utf8.EncodeRune(rb[:], lo)

				// 跳过已处理的字节
				rb = shiftNRuneBytes(rb, off)

				idxc := rb[0]
				for i, c := range []byte(n.indices) {
					// 小写匹配
					if c == idxc {
						// 必须使用递归方法，因为大写字节和小写字节可能同时存在
						if out := n.children[i].findCaseInsensitivePathRec(
							path, ciPath, rb, fixTrailingSlash,
						); out != nil {
							return out
						}
						break
					}
				}

				// 如果没有找到匹配，则对大写字母执行相同的操作
				if up := unicode.ToUpper(rv); up != lo {
					utf8.EncodeRune(rb[:], up)
					rb = shiftNRuneBytes(rb, off)

					idxC := rb[0]
					for i, c := range []byte(n.indices) {
						// 大写匹配
						if c == idxC {
							// 继续处理子节点
							n = n.children[i]
							npLen = len(n.path)
							continue walk
						}
					}
				}
			}

			// 没有找到。尝试通过添加/删除尾随斜杠来修复路径
			if fixTrailingSlash && path == constants.PathSeparatorStr && n.handlers != nil {
				return ciPath
			}
			return nil
		}

		n = n.children[0]
		switch n.nType { //nolint:exhaustive
		case paramNode:
			// 查找 paramNode 结束（要么是 constants.PathSeparator 要么是路径结束）
			end := 0
			for end < len(path) && path[end] != constants.PathSeparator {
				end++
			}

			// 将 paramNode 值添加到不区分大小写的路径中
			ciPath = append(ciPath, path[:end]...)

			// 我们需要更深入！
			if end < len(path) {
				if len(n.children) > 0 {
					// 继续处理子节点
					n = n.children[0]
					npLen = len(n.path)
					path = path[end:]
					continue
				}

				// ...但我们不能
				if fixTrailingSlash && len(path) == end+1 {
					return ciPath
				}
				return nil
			}

			if n.handlers != nil {
				return ciPath
			}

			if fixTrailingSlash && len(n.children) == 1 {
				// 没有找到处理函数。检查此路径 + 尾随斜杠是否存在处理函数
				n = n.children[0]
				if n.path == constants.PathSeparatorStr && n.handlers != nil {
					return append(ciPath, constants.PathSeparator)
				}
			}

			return nil

		case wildcardNode:
			return append(ciPath, path...)

		default:
			panic("无效的节点类型")
		}
	}
	// 没有找到匹配的处理函数。
	// 尝试通过添加或删除尾随斜杠来修复路径
	if fixTrailingSlash {
		// 如果当前路径是单个斜杠，则直接返回当前路径
		if path == constants.PathSeparatorStr {
			return ciPath // 返回不区分大小写的路径
		}
		// 检查路径长度与节点路径长度的关系
		// 如果路径长度加 1 等于节点路径长度，并且节点路径的下一个字符是 constants.PathSeparator
		// 并且当前路径（去掉前导斜杠）与节点路径（去掉前导斜杠）相等
		// 且当前节点注册了处理函数
		if len(path)+1 == npLen && n.path[len(path)] == constants.PathSeparator &&
			strings.EqualFold(path[1:], n.path[1:len(path)]) && n.handlers != nil {
			// 将节点路径添加到不区分大小写的路径中，返回修复后的路径
			return append(ciPath, n.path...)
		}
	}
	// 如果没有找到匹配的处理函数并且没有进行路径修复，则返回 nil
	return nil

}
