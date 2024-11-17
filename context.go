/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 22:36:15
 * @FilePath: \gosh\context.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/gosh/constants" // 引入常量定义
	"github.com/kamalyes/gosh/errorsx"   // 引入error定义
)

// ContextInterface 定义了 Context 结构体暴露的方法
type ContextInterface interface {
	// 路径和查询参数处理
	PathValue(key string) string          // 根据键获取路径参数值
	PathParam(key string) (string, bool)  // 获取路径参数，同时返回是否存在
	AllPathValues() []Param               // 获取所有路径参数
	QueryValue(key string) string         // 获取查询参数值
	QueryParam(key string) (string, bool) // 获取查询参数，同时返回是否存在
	AllQueryValues() url.Values           // 获取所有查询参数
	FormValue(key string) string          // 获取表单参数值
	AllFormValues() url.Values            // 获取所有表单参数

	// 上下文处理
	SetContextValue(key, value any) // 设置上下文中的值
	GetContextValue(key any) any    // 获取上下文中的值
	Deadline() (time.Time, bool)    // Deadline 返回请求的截止时间和一个布尔值，表示是否存在截止时间。如果请求没有上下文，则返回零值和 false。
	Done() <-chan struct{}          // Done 返回一个通道，当请求的上下文完成时，该通道会关闭。如果请求没有上下文，则返回 nil（表示永远不会完成）。
	Err() error                     // Err 返回请求上下文的错误信息。如果请求没有上下文，则返回 nil。
	Copy() *Context                 // 复制上下文

	// 响应处理
	WriteString(status int, data string) error        // 返回字符串响应
	WriteJSONResponse(status int, data any) error     // 返回 JSON 响应
	WriteNoContent() error                            // 无内容响应
	WriteRedirect(code int, url string) error         // 重定向响应
	WriteJSONWithStatus(status int, data any) error   // 设置 JSON 响应并指定状态码
	AbortWithStatus(code int)                         // 中止请求并返回状态码
	AbortWithStatusJSON(code int, jsonObj any)        // 中止请求并返回 JSON 状态错误信息
	AbortWithStatusText(code int, message string)     // 中止请求 Text 的处理方法
	AbortWithStatusHTML(code int, htmlContent string) // 中止请求 HTML 状态处理方法
	AbortWithError(code int, err error) error         // 中止请求并返回错误信息

	// 文件处理
	ServeFile(filePath string) error                                                    // 提供指定路径的文件
	FileFromFS(filePath string, fs http.FileSystem) error                               // 从文件系统提供文件
	SaveFile(fileHeader *multipart.FileHeader, savePath string, perm os.FileMode) error // 保存上传的文件

	// 请求处理
	Method() string           // 获取请求方法
	Path() string             // 获取请求路径
	Header(key string) string // 获取请求头
	AllHeaders() http.Header  // 获取所有请求头
	ClientIP() string         // 获取客户端 IP 地址
	UserAgent() string        // 获取请求的 User-Agent

	// Cookie 和请求体处理
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) // 设置 Cookie
	Cookie(name string) (string, error)                                                   // 获取 Cookie
	Body() ([]byte, error)                                                                // 获取请求体内容
	JSONParseBody(obj any) error                                                          // 获取 Json请求体
	IsMethod(method string) bool                                                          // 检查请求方法是否为指定的方法
	FullRequestPath() string                                                              // 获取请求的完整路径
	GetURLParam(key string) string                                                        // 获取 URL 参数
	MultipartForm() (*multipart.Form, error)                                              // 解析后的多部分表单
}

// Context 是处理 HTTP 请求的核心结构体
type Context struct {
	Request        *http.Request       // 当前 HTTP 请求
	ResponseWriter http.ResponseWriter // HTTP 响应写入器
	Status         int                 // HTTP 处理状态码
	Error          error               // 错误信息
	broke          bool                // 请求是否被中止的标志
	index          int8                // 当前处理程序的索引
	fullPath       string              // 请求的完整路径
	Engine         *Engine             // 关联的引擎（例如框架）
	params         *Params             // 路径参数
	skippedNodes   *[]skippedNode      // 被跳过的节点
	queryCache     url.Values          // 查询参数缓存
	formCache      url.Values          // 表单参数缓存
	handlers       HandlersChain       // 处理程序链
}

// 实现 ContextInterface
var _ ContextInterface = (*Context)(nil) // 确保 Context 实现了 ContextInterface

// 重置上下文状态
func (ctx *Context) reset() {
	ctx.Request = &http.Request{}               // 重置请求
	ctx.Status = http.StatusOK                  // 重置状态为 200 OK
	ctx.Error = nil                             // 清空错误信息
	ctx.index = -1                              // 处理程序索引重置
	ctx.broke = false                           // 请求未被中止
	ctx.fullPath = ""                           // 清空完整路径
	ctx.queryCache = nil                        // 清空查询参数缓存
	ctx.formCache = nil                         // 清空表单参数缓存
	*ctx.params = (*ctx.params)[:0]             // 清空路径参数
	*ctx.skippedNodes = (*ctx.skippedNodes)[:0] // 清空被跳过的节点
}

// 获取引擎配置
func (ctx *Context) EngineConfig() Config {
	return ctx.Engine.Config // 返回与当前上下文相关的引擎配置
}

// 返回路由完整路径
func (ctx *Context) FullPath() string {
	return ctx.fullPath // 返回请求的完整路径
}

// 停止当前请求的处理
func (ctx *Context) Abort() *Context {
	ctx.broke = true // 将请求标记为已中止
	return ctx
}

// 判断请求是否已被中止
func (ctx *Context) IsAborted() bool {
	return ctx.broke // 返回请求是否已被中止
}

// AbortWithStatus 调用 `Abort()` 并使用指定的状态码写入响应头。
func (c *Context) AbortWithStatus(code int) {
	c.Status = code                    // 直接赋值给 Status 字段
	c.ResponseWriter.WriteHeader(code) // 写入响应头
	c.Abort()                          // 中止请求处理
}

// AbortWithStatusJSON 调用 `Abort()`，然后内部调用 `JSON`。
func (c *Context) AbortWithStatusJSON(code int, jsonObj any) {
	c.Abort()                          // 中止请求处理
	c.WriteJSONResponse(code, jsonObj) // 写入 JSON 响应
}

// AbortWithStatusText 调用 `Abort()` 并使用指定的状态码和文本写入响应头。
func (c *Context) AbortWithStatusText(code int, message string) {
	c.Status = code                              // 设置状态码
	c.setContentType(constants.ContentTypePlain) // 设置 Content-Type 为文本
	c.ResponseWriter.WriteHeader(code)           // 写入响应头
	c.ResponseWriter.Write([]byte(message))      // 写入文本消息
}

// AbortWithStatusHTML 调用 `Abort()` 并使用指定的状态码和 HTML 内容写入响应头。
func (c *Context) AbortWithStatusHTML(code int, htmlContent string) {
	c.Abort()                                   // 中止请求处理
	c.Status = code                             // 设置状态码
	c.setContentType(constants.ContentTypeHtml) // 设置 Content-Type 为 HTML
	c.ResponseWriter.WriteHeader(code)          // 写入响应头
	c.ResponseWriter.Write([]byte(htmlContent)) // 写入 HTML 内容
}

// AbortWithError 调用 `AbortWithStatus()` 和 `Error()`。
func (c *Context) AbortWithError(code int, err error) error {
	c.AbortWithStatus(code) // 中止请求处理并写入状态码
	return err              // 记录错误信息
}

// 执行下一个处理程序
func (ctx *Context) Next() {
	ctx.index++                               // 增加处理程序索引
	for ctx.index < int8(len(ctx.handlers)) { // 遍历处理程序链
		ctx.handlers[ctx.index](ctx) // 执行当前处理程序
		ctx.index++                  // 移动到下一个处理程序
	}
}

// SetContextValue 设置上下文中的值
func (ctx *Context) SetContextValue(key, value any) {
	if key != nil {
		ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), key, value))
	}
}

// 从上下文中获取值
func (ctx *Context) GetContextValue(key any) any {
	if key != nil {
		return ctx.Request.Context().Value(key) // 从上下文中获取指定键的值
	}
	return nil
}

// hasRequestContext 判断请求是否有上下文
func (c *Context) hasRequestContext() bool {
	hasFallback := c.Engine != nil
	hasRequestContext := c.Request != nil && c.Request.Context() != nil
	return hasFallback && hasRequestContext
}

// Deadline 返回请求的截止时间
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	if !c.hasRequestContext() {
		return
	}
	return c.Request.Context().Deadline()
}

// Done 返回请求的完成信号
func (c *Context) Done() <-chan struct{} {
	if !c.hasRequestContext() {
		return nil
	}
	return c.Request.Context().Done()
}

// Err 返回请求的错误信息
func (c *Context) Err() error {
	if !c.hasRequestContext() {
		return nil
	}
	return c.Request.Context().Err()
}

// Copy 返回当前上下文的副本，可以安全地在请求范围外使用。
// 当需要将上下文传递给 goroutine 时，必须使用此方法。
func (c *Context) Copy() *Context {
	cp := &Context{
		Request:    c.Request,        // 复制请求
		Status:     c.Status,         // 复制状态
		Error:      c.Error,          // 复制错误信息
		broke:      c.broke,          // 复制中止标志
		index:      c.index,          // 复制处理程序索引
		fullPath:   c.fullPath,       // 复制完整路径
		Engine:     c.Engine,         // 复制引擎
		params:     c.params,         // 复制路径参数
		queryCache: make(url.Values), // 创建新的查询参数缓存
		formCache:  make(url.Values), // 创建新的表单参数缓存
		handlers:   nil,              // 清空处理程序链
	}

	// 复制查询参数缓存
	for k, v := range c.queryCache {
		cp.queryCache[k] = v
	}

	// 复制表单参数缓存
	for k, v := range c.formCache {
		cp.formCache[k] = v
	}

	// 处理跳过的节点（如果需要）
	if c.skippedNodes != nil {
		cp.skippedNodes = new([]skippedNode)
		*cp.skippedNodes = make([]skippedNode, len(*c.skippedNodes))
		copy(*cp.skippedNodes, *c.skippedNodes)
	}

	// 复制其他字段（如需要）
	cp.handlers = append([]HandlerFunc{}, c.handlers...) // 复制处理程序链

	return cp // 返回副本
}

// 路径参数处理
func (ctx *Context) PathValue(key string) string {
	return ctx.params.ByName(key) // 根据键获取路径参数值
}

// 获取路径参数，同时返回是否存在
func (ctx *Context) PathParam(key string) (string, bool) {
	return ctx.params.Get(key) // 获取路径参数并返回存在性
}

// 获取所有路径参数
func (ctx *Context) AllPathValues() []Param {
	return *ctx.params // 返回所有路径参数
}

// 初始化查询参数缓存
func (ctx *Context) initQueryCache() {
	if ctx.queryCache == nil { // 如果查询缓存尚未初始化
		ctx.queryCache = ctx.Request.URL.Query() // 从请求中读取查询参数
	}
}

// 查询参数处理
func (ctx *Context) QueryValue(key string) string {
	ctx.initQueryCache()           // 确保查询缓存已初始化
	return ctx.queryCache.Get(key) // 返回指定键的查询参数值
}

// 获取查询参数，同时返回是否存在
func (ctx *Context) QueryParam(key string) (string, bool) {
	ctx.initQueryCache()              // 确保查询缓存已初始化
	values, ok := ctx.queryCache[key] // 获取查询参数
	if !ok {
		return "", false // 如果不存在，则返回空值和 false
	}
	return values[0], true // 返回第一个匹配的值和存在性
}

// 获取所有查询参数
func (ctx *Context) AllQueryValues() url.Values {
	ctx.initQueryCache()  // 确保查询缓存已初始化
	return ctx.queryCache // 返回所有查询参数
}

// 表单参数处理
func (ctx *Context) InitFormCache() error {
	if ctx.formCache != nil { // 如果表单缓存已初始化
		return nil
	}
	ctx.formCache = make(url.Values) // 创建新的表单参数缓存

	// 解析请求的表单数据
	if err := ctx.Request.ParseMultipartForm(ctx.Engine.Config.MaxMultipartMemory); err != nil && !errors.Is(err, errorsx.ErrNotMultipart) {
		return err // 如果解析失败，返回错误
	}
	ctx.formCache = ctx.Request.PostForm // 将解析的表单数据存入缓存
	return nil
}

// 获取表单参数值
func (ctx *Context) FormValue(key string) string {
	if err := ctx.InitFormCache(); err != nil { // 确保表单缓存已初始化
		return ""
	}
	return ctx.formCache.Get(key) // 返回表单参数值
}

// 获取所有表单参数
func (ctx *Context) AllFormValues() url.Values {
	if err := ctx.InitFormCache(); err != nil { // 确保表单缓存已初始化
		return nil
	}
	return ctx.formCache // 返回所有表单参数
}

// MultipartForm 是解析后的多部分表单，包括文件上传。
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(c.Engine.Config.MaxMultipartMemory)
	return c.Request.MultipartForm, err
}

// 文件处理
func (ctx *Context) FormFile(name string) (*multipart.FileHeader, error) {
	if err := ctx.InitFormCache(); err != nil { // 确保表单缓存已初始化
		return nil, err
	}
	f, fileHeader, err := ctx.Request.FormFile(name) // 获取指定名称的文件
	if err != nil {
		return nil, err // 返回错误信息
	}
	defer f.Close()

	return fileHeader, nil // 返回文件头部信息
}

// 提供文件的公共逻辑
func (ctx *Context) serveFileResponse(file io.Reader, fileInfo os.FileInfo) error {
	contentType := constants.ContentTypeOctet // 默认内容类型
	if ext := filepath.Ext(fileInfo.Name()); ext != "" {
		contentType = mime.TypeByExtension(ext) // 根据文件扩展名设置内容类型
	}
	// 设置 Content-Type 为文本
	ctx.setContentType(contentType)
	// 设置 Content-Length 头部
	ctx.setHeaderContentLength(convert.MustString(fileInfo.Size()))

	// 写入响应头
	ctx.ResponseWriter.WriteHeader(http.StatusOK) // 设置响应状态为 200 OK

	// 将文件内容写入响应
	_, err := io.Copy(ctx.ResponseWriter, file) // 将文件内容拷贝到响应写入器
	return err
}

// ServeFile 提供文件给客户端并设置适当的 Content-Type 头部
func (ctx *Context) ServeFile(filePath string) error {
	file, err := os.Open(filePath) // 打开指定路径的文件
	if err != nil {
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound) // 如果文件未找到，返回 404
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat() // 获取文件信息
	if err != nil {
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError) // 如果获取文件信息失败，返回 500
		return err
	}

	return ctx.serveFileResponse(file, fileInfo) // 调用公共方法返回文件内容
}

// FileFromFS 从文件系统中提供文件
func (ctx *Context) FileFromFS(filePath string, fs http.FileSystem) error {
	file, err := fs.Open(filePath) // 从文件系统打开文件
	if err != nil {
		ctx.ResponseWriter.WriteHeader(http.StatusNotFound) // 如果文件未找到，返回 404
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat() // 获取文件信息
	if err != nil {
		ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError) // 如果获取文件信息失败，返回 500
		return err
	}

	return ctx.serveFileResponse(file, fileInfo) // 调用公共方法返回文件内容
}

// SaveFile 保存上传的文件到指定路径
func (ctx *Context) SaveFile(fileHeader *multipart.FileHeader, savePath string, perm os.FileMode) error {
	f, err := fileHeader.Open() // 打开上传的文件
	if err != nil {
		return err // 返回错误信息
	}
	defer f.Close()

	if err := os.MkdirAll(filepath.Dir(savePath), perm); err != nil { // 创建保存路径
		return err
	}

	out, err := os.Create(savePath) // 创建目标文件
	if err != nil {
		return err // 返回错误信息
	}
	defer out.Close()

	_, err = io.Copy(out, f) // 将上传的文件内容拷贝到目标文件
	return err
}

// H is a shortcut for map[string]any
type H map[string]any

// WrapF 是一个辅助函数，用于包装 http.HandlerFunc 并返回一个gosh中间件。
func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		f(c.ResponseWriter, c.Request) // 调用传入的 http.HandlerFunc，处理响应和请求
		return nil                     // 返回 nil，表示没有错误
	}
}

// WrapH 是一个辅助函数，用于包装 http.Handler 并返回一个gosh中间件。
func WrapH(h http.Handler) HandlerFunc {
	return func(c *Context) error {
		h.ServeHTTP(c.ResponseWriter, c.Request) // 调用传入的 http.Handler，处理响应和请求
		return nil                               // 返回 nil，表示没有错误
	}
}

// 响应处理
func (ctx *Context) WriteString(status int, data string) error {
	ctx.setContentType(constants.ContentTypePlain)                      // 设置 Content-Type 为文本
	ctx.ResponseWriter.WriteHeader(status)                              // 写入响应状态码
	_, err := ctx.ResponseWriter.Write(convert.StringToSliceByte(data)) // 将字符串写入响应
	return err
}

// 返回 JSON 响应
func (ctx *Context) WriteJSONResponse(status int, data any) error {
	buf, err := json.Marshal(data) // 将数据序列化为 JSON
	if err != nil {
		return err // 返回错误信息
	}
	ctx.setContentType(constants.ContentTypeJSON) // 设置 Content-Type 为 JSON
	// ctx.ResponseWriter.Header().Set(constants.HeaderContentLength, strconv.Itoa(len(buf))) // 设置 Content-Length
	ctx.ResponseWriter.WriteHeader(status) // 写入响应状态码
	_, err = ctx.ResponseWriter.Write(buf) // 将 JSON 数据写入响应
	return err
}

// 无内容响应
func (ctx *Context) WriteNoContent() error {
	ctx.Status = http.StatusNoContent                    // 状态设为 204 No Content
	ctx.ResponseWriter.WriteHeader(http.StatusNoContent) // 写入响应状态
	return nil
}

// 重定向响应
func (ctx *Context) WriteRedirect(code int, url string) error {
	if code < 300 || code > 308 { // 检查重定向状态码
		return errorsx.ErrInvalidRedirectCode // 返回无效状态码错误
	}
	ctx.Status = code                                              // 设置状态码
	ctx.ResponseWriter.Header().Set(constants.HeaderLocation, url) // 设置 Location 头部
	ctx.ResponseWriter.WriteHeader(code)                           // 写入响应状态
	return nil
}

// Json解析请求Body数据
func (ctx *Context) JSONParseBody(obj any) error {
	body, err := io.ReadAll(ctx.Request.Body) // 读取请求体
	if err != nil {
		return err // 返回错误信息
	}
	return json.Unmarshal(body, obj) // 将 JSON 数据解析到目标对象中
}

// 获取请求方法
func (ctx *Context) Method() string {
	return ctx.Request.Method // 返回请求的方法
}

// 获取请求路径
func (ctx *Context) Path() string {
	return ctx.Request.URL.Path // 返回请求的路径
}

// 获取请求头
func (ctx *Context) Header(key string) string {
	return ctx.Request.Header.Get(key) // 根据键获取请求头的值
}

// 获取所有请求头
func (ctx *Context) AllHeaders() http.Header {
	return ctx.Request.Header // 返回所有请求头
}

// 获取请求的客户端 IP 地址
func (ctx *Context) ClientIP() string {
	ip := ctx.Request.Header.Get("X-Forwarded-For") // 尝试从 X-Forwarded-For 头获取 IP
	if ip == "" {
		ip, _, _ = net.SplitHostPort(ctx.Request.RemoteAddr) // 如果没有，则从 RemoteAddr 获取
	}
	return ip // 返回客户端 IP 地址
}

// 获取请求的 User-Agent
func (ctx *Context) UserAgent() string {
	return ctx.Request.UserAgent() // 返回请求的 User-Agent
}

// 设置 Content-Type 头部
func (ctx *Context) setContentType(contentType string) {
	ctx.ResponseWriter.Header().Set(constants.HeaderContentType, contentType) // 设置响应的 Content-Type
}

// 设置 HeaderContentLength 头部
func (ctx *Context) setHeaderContentLength(len string) {
	ctx.ResponseWriter.Header().Set(constants.HeaderContentLength, len) // 设置响应的 Content-Length
}

// SetStatus 设置响应的状态码
func (ctx *Context) SetStatus(status int) {
	ctx.Status = status
}

// SetHeader 设置响应头的指定键值
func (ctx *Context) SetHeader(key, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}

// SetCookie 设置 Cookie
func (ctx *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	http.SetCookie(ctx.ResponseWriter, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// Cookie 返回请求中提供的指定名称的 cookie，
// 如果未找到则返回 ErrNoCookie。
// 返回的 cookie 值会被解码（unescaped）。
// 如果有多个 cookie 匹配给定的名称，则只返回其中一个。
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err // 如果未找到 cookie，返回错误
	}
	val, _ := url.QueryUnescape(cookie.Value) // 解码 cookie 值
	return val, nil                           // 返回解码后的 cookie 值
}

// Body 获取请求体内容
func (ctx *Context) Body() ([]byte, error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 重新设置请求体
	return body, nil
}

// IsMethod 检查请求方法是否为指定的方法
func (ctx *Context) IsMethod(method string) bool {
	return ctx.Request.Method == method
}

// FullRequestPath 返回请求的完整路径（包括查询参数）
func (ctx *Context) FullRequestPath() string {
	return ctx.Request.URL.String()
}

// WriteJSONWithStatus 设置 JSON 响应并指定状态码
func (ctx *Context) WriteJSONWithStatus(status int, data any) error {
	ctx.SetStatus(status)
	return ctx.WriteJSONResponse(status, data)
}

// GetURLParam 获取 URL 参数
func (ctx *Context) GetURLParam(key string) string {
	return ctx.PathValue(key)
}
