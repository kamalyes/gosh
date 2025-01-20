/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 13:26:00
 * @FilePath: \gosh\engine.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"

	translator "github.com/go-playground/universal-translator"
	goconfig "github.com/kamalyes/go-config"
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/random"
	"github.com/kamalyes/gosh/constants"
	"github.com/kamalyes/gosh/errorsx"
)

// 常量定义
const (
	defaultMaxMultipartMemory = 32 << 20 // 默认32 MB
)

// Config 引擎参数配置
type Config struct {
	MaxMultipartMemory     int64                  // 允许的请求Body大小(默认32 << 20 = 32MB)
	Recovery               bool                   // 自动恢复panic，防止进程退出
	HandleMethodNotAllowed bool                   // 是否处理 405 错误（可以减少路由匹配时间），以 404 错误返回
	BeforeHandler          CallbackHandler        // 前置回调处理器，总是会在其它处理器执行之前执行
	ErrorHandler           CallbackHandler        // 错误回调处理器
	AfterHandler           CallbackHandler        // 后置回调处理器，总是会在其它处理器全部执行完之后执行
	AppBanner              *BannerConfig          // Banner配置
	AppName                string                 // 应用名称
	Zap                    *Logger                // 日志
	Trans                  translator.Translator  // Trans 全局validate翻译器
	KmSingleConfig         *goconfig.SingleConfig // 私有配置
}

// HandlerFunc 路由处理器函数类型
type HandlerFunc func(*Context) error

// HandlersChain 处理器链
type HandlersChain []HandlerFunc

// CallbackHandler 回调处理器
type CallbackHandler func(*Context)

// Engine 引擎
type Engine struct {
	RouterGroup
	Config      Config       // 引擎配置
	maxParams   int          // 最大参数数量
	maxSections int          // 最大路径段数量
	contextPool sync.Pool    // 上下文池
	trees       methodTrees  // 路由树
	routes      []*RouteInfo // 存储路由，使用 RouteInfo 结构体
}

// NewEngine 新建引擎实例
func NewEngine(config ...Config) *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			handlers: nil,
			basePath: constants.PathSeparatorStr,
			root:     true,
		},
		trees: make(methodTrees, 0, 9), // 初始化路由树
	}
	engine.RouterGroup.Engine = engine
	engine.Config = setDefaultConfig()

	if len(config) > 0 {
		engine.Config = mergeDefaultConfig(engine.Config, config[0])
	}

	// 初始化上下文池
	engine.contextPool.New = func() any {
		return engine.allocateContext(engine.maxParams)
	}

	return engine
}

// Run 启动HTTP服务
func (engine *Engine) Run(addr ...string) error {
	resolveAddress := resolveAddress(addr)
	engine.Config.AppBanner.Print()
	log.Printf("Starting server at %s", resolveAddress)
	return http.ListenAndServe(resolveAddress, engine)
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			return constants.WildcardSymbolStr + port
		}
		port, err := random.GenerateAvailablePort()
		if err != nil {
			panic(err)
		}
		return fmt.Sprintf("%s%d", constants.WildcardSymbolStr, port)
	case 1:
		return addr[0]
	default:
		panic("resolveAddress too many parameters")
	}
}

// 检查路由是否存在
func (e *Engine) routeExists(method, path string) (bool, *RouteInfo) {
	for _, routeInfo := range e.routes {
		if routeInfo.Method == method && routeInfo.Path == path {
			return true, routeInfo
		}
	}
	return false, nil // 返回 nil 表示未找到
}

// GetAllRoutes 返回所有的路由信息
func (e *Engine) GetAllRoutes() []*RouteInfo {
	return e.routes
}

// allocateContext 分配上下文
func (engine *Engine) allocateContext(maxParams int) *Context {
	v := make(Params, 0, maxParams)                            // 创建参数
	skippedNodes := make([]skippedNode, 0, engine.maxSections) // 创建跳过的节点

	return &Context{
		Engine:       engine,
		params:       &v,
		skippedNodes: &skippedNodes,
	}
}

// String 方法返回处理程序链的字符串表示
func (hc HandlersChain) String() string {
	var names []string
	for _, handler := range hc {
		names = append(names, getHandlerName(handler)) // 假设有一个函数可以获取处理程序的名称
	}
	return "[" + strings.Join(names, ", ") + "]"
}

// getHandlerName 返回处理程序函数的名称
func getHandlerName(handler HandlerFunc) string {
	// 获取函数的调用信息
	funcPtr := runtime.FuncForPC(reflect.ValueOf(handler).Pointer())
	if funcPtr == nil {
		return "unknown"
	}

	// 获取函数名称
	funcName := funcPtr.Name()
	// 提取包名和函数名
	parts := strings.Split(funcName, ".")
	if len(parts) > 1 {
		return parts[len(parts)-1] // 返回函数名
	}
	return funcName
}

// addRoute 添加路由
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	validateRoute(method, path, handlers)

	root := engine.trees.get(method) // 获取指定方法的根节点
	if root == nil {
		root = new(Node) // 创建新的根节点
		root.fullPath = constants.PathSeparatorStr
		engine.trees = append(engine.trees, methodTree{method: method, root: root}) // 添加到树中
	}
	root.addRoute(path, handlers) // 添加路由处理器

	engine.updateMaxParamsAndSections(path)
	engine.updateRoutes(method, path, handlers)
}

// validateRoute 验证路由合法性
func validateRoute(method, path string, handlers HandlersChain) {
	if path[0] != constants.PathSeparator {
		log.Fatalln(errorsx.ErrPathMustStartWithSlash)
	}
	if method == "" {
		log.Fatalln(errorsx.ErrMethodCannotBeEmpty)
	}
	if len(handlers) == 0 {
		log.Fatalln(errorsx.ErrMustHaveAtLeastOneHandler)
	}
}

// updateRoutes 更新路由表信息
func (engine *Engine) updateRoutes(method, path string, handler HandlersChain) {
	// 创建一个新的 RouteInfo
	routeInfo := &RouteInfo{
		Method:  method,
		Path:    path,
		Handler: handler,
	}

	// 检查是否已经存在相同的路由
	for i, existingRoute := range engine.routes {
		if existingRoute.Method == method && existingRoute.Path == path {
			// 更新现有的路由
			engine.routes[i] = routeInfo
			return
		}
	}

	// 添加新路由
	engine.routes = append(engine.routes, routeInfo)
}

// updateMaxParamsAndSections 更新最大参数数量和路径段数量
func (engine *Engine) updateMaxParamsAndSections(path string) {
	if paramsCount := mathx.CountPathSegments(path, constants.PathSeparatorStr, constants.PathParamPrefixStr); paramsCount > engine.maxParams {
		engine.maxParams = paramsCount
	}

	if sectionsCount := mathx.CountPathSegments(path, constants.PathSeparatorStr); sectionsCount > engine.maxSections {
		engine.maxSections = sectionsCount
	}
}

// ServeHTTP 处理HTTP请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := engine.prepareContext(w, req) // 从池中获取并准备上下文

	// 处理panic
	if engine.Config.Recovery {
		defer engine.recoverFromPanic(ctx)
	}

	engine.handleRequest(ctx)   // 处理请求
	engine.contextPool.Put(ctx) // 将上下文放回池中
}

// prepareContext 准备上下文
func (engine *Engine) prepareContext(w http.ResponseWriter, req *http.Request) *Context {
	ctx := engine.contextPool.Get().(*Context) // 从池中获取上下文
	ctx.reset()                                // 重置上下文
	// 初始化上下文
	ctx.Request = req
	ctx.ResponseWriter = w
	return ctx
}

// recoverFromPanic 处理panic
func (engine *Engine) recoverFromPanic(ctx *Context) {
	if err := recover(); err != nil {
		// 创建自定义错误，附带 ErrorTypePrivate 表示这是一个私有错误
		customErr := &errorsx.CustomError{
			Err:       fmt.Errorf("%v", err),
			ErrorType: errorsx.ErrorTypePrivate,
		}
		handleError(ctx, engine, customErr, http.StatusInternalServerError)
	}
}

// handleRequest 处理请求的核心逻辑
func (engine *Engine) handleRequest(ctx *Context) {
	method := ctx.Request.Method // 获取请求方法
	url := ctx.Request.URL.Path  // 获取请求路径
	node, found := engine.findNode(method, url, ctx)

	// 如果找不到 返回错误
	if !found {
		engine.handleNotFoundOrMethodNotAllowed(ctx, method, url)
		return
	}
	// OK 即正常逻辑
	if node.handlers != nil {
		engine.executeHandlers(node, ctx)
	}
}

// findNode 查找路由节点
func (engine *Engine) findNode(method, url string, ctx *Context) (nodeValue, bool) {
	for _, tree := range engine.trees {
		if tree.method != method {
			continue
		}
		root := tree.root
		node := root.getValue(url, ctx.params, ctx.skippedNodes) // 查找路由节点

		if node.params != nil {
			ctx.params = node.params
		}
		return node, true
	}
	return nodeValue{}, false
}

// executeHandlers 执行路由处理器
func (engine *Engine) executeHandlers(node nodeValue, ctx *Context) {
	ctx.fullPath = node.fullPath

	if engine.Config.BeforeHandler != nil {
		engine.Config.BeforeHandler(ctx)
	}

	defer func() {
		if engine.Config.AfterHandler != nil {
			engine.Config.AfterHandler(ctx)
		}
	}()

	for _, handler := range node.handlers {
		if ctx.broke {
			break
		}
		// 执行处理器并检查错误
		if err := handler(ctx); err != nil {
			engine.handleError(ctx, err) // 将 err 传递给 handleError 方法
			return
		}
	}
}

// handleError 封装错误处理
func (engine *Engine) handleError(ctx *Context, err error) {
	var customErr *errorsx.CustomError
	// 尝试将 err 转换为 *errorsx.CustomError
	if errors.As(err, &customErr) {
		engine.processError(ctx, customErr, http.StatusInternalServerError)
		return
	}

	// 创建新的 CustomError
	customErr = &errorsx.CustomError{
		Err:       err,
		ErrorType: errorsx.ErrorTypePrivate, // 指定错误类型
	}
	engine.processError(ctx, customErr, http.StatusInternalServerError)
}

// processError 处理错误并执行错误处理器
func (engine *Engine) processError(ctx *Context, err *errorsx.CustomError, status int) {
	ctx.broke = true // 标记上下文为中断状态
	ctx.Status = status
	ctx.Error = err

	if engine.Config.ErrorHandler != nil {
		engine.Config.ErrorHandler(ctx) // 调用错误处理器
		return
	}
	ctx.ResponseWriter.WriteHeader(ctx.Status)

	// 直接使用自定义错误的字符串表示
	if _, errWrite := ctx.ResponseWriter.Write(convert.StringToSliceByte(err.Error())); errWrite != nil {
		log.Println("写入响应时出错:", err)
	}
}

// handleNotFoundOrMethodNotAllowed 处理404或405错误
func (engine *Engine) handleNotFoundOrMethodNotAllowed(ctx *Context, method, url string) {
	if !engine.Config.HandleMethodNotAllowed || !engine.isMethodAllowed(method, url) {
		handleError(ctx, engine, errorsx.ErrNotFound, http.StatusNotFound)
		return
	}
	handleError(ctx, engine, errorsx.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
}

// isMethodAllowed 检查方法是否被允许
func (engine *Engine) isMethodAllowed(method, url string) bool {
	for _, tree := range engine.trees {
		if tree.method != method {
			node := tree.root.getValue(url, nil, nil)
			if node.handlers != nil {
				return true
			}
		}
	}
	return false
}

// handleError 处理错误并执行错误处理器
func handleError(ctx *Context, engine *Engine, err *errorsx.CustomError, status int) error {
	ctx.broke = true // 标记上下文为中断状态
	ctx.Status = status
	ctx.Error = err

	if engine.Config.ErrorHandler != nil {
		engine.Config.ErrorHandler(ctx) // 调用错误处理器
		return nil
	}
	ctx.ResponseWriter.WriteHeader(ctx.Status)

	// 直接使用自定义错误的字符串表示
	if _, errWrite := ctx.ResponseWriter.Write(convert.StringToSliceByte(err.Error())); errWrite != nil {
		log.Println("写入响应时出错:", err)
		return errWrite
	}
	return nil
}
