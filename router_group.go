/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 12:05:15
 * @FilePath: \gosh\router_group.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"fmt"
	"net/http"

	"github.com/kamalyes/go-toolbox/pkg/osx"
)

// RouterGroup 路由组
type RouterGroup struct {
	handlers HandlersChain // 路由组中注册的处理程序链
	basePath string        // 路由组的基础路径
	Engine   *Engine       // 引擎实例
	root     bool          // 是否为根路由组
	noRoute  HandlersChain // 没有匹配路由时的处理程序
}

// RouteInfo 表示请求路由的规范，包括请求方法、路径及其处理函数。
type RouteInfo struct {
	Method  string        // 请求方法，例如 GET、POST 等
	Path    string        // 请求路径
	Handler HandlersChain // 实际的处理函数
}

// NoRoute 注册没有匹配路由时的处理程序
func (group *RouterGroup) NoRoute(handlers ...HandlerFunc) {
	group.noRoute = handlers
}

// Use 使用中间件
func (group *RouterGroup) Use(handlers ...HandlerFunc) {
	group.handlers = append(group.handlers, handlers...)
}

// Group 注册路由组
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		Engine:   group.Engine,
	}
}

// handle 处理路由注册
func (group *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) error {
	absolutePath := group.calculateAbsolutePath(relativePath)

	// 检查是否已经存在相同的路径和方法的路由
	if exists, existingRoute := group.Engine.routeExists(httpMethod, absolutePath); exists {
		panic(fmt.Sprintf("route already exists: %s %s with handlers: %s", httpMethod, absolutePath, existingRoute.Handler.String()))
	}

	// 合并处理程序链
	handlers = group.combineHandlers(handlers)
	group.Engine.addRoute(httpMethod, absolutePath, handlers)
	return nil
}

// Handle 注册自定义方法的路由
func (group *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) error {
	return group.handle(httpMethod, relativePath, handlers)
}

// 各种 HTTP 方法的路由注册
func (group *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) error {
	return group.handle(http.MethodPost, relativePath, handlers)
}

func (group *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) error {
	return group.handle(http.MethodGet, relativePath, handlers)
}

func (group *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) error {
	return group.handle(http.MethodDelete, relativePath, handlers)
}

func (group *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) error {
	return group.handle(http.MethodPatch, relativePath, handlers)
}

func (group *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) error {
	return group.handle(http.MethodPut, relativePath, handlers)
}

func (group *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) error {
	return group.handle(http.MethodOptions, relativePath, handlers)
}

func (group *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) error {
	return group.handle(http.MethodHead, relativePath, handlers)
}

// Match 根据方法数组和路径注册路由
func (group *RouterGroup) Match(methods []string, relativePath string, handlers ...HandlerFunc) error {
	for _, method := range methods {
		if err := group.handle(method, relativePath, handlers); err != nil {
			return err // 返回第一个错误
		}
	}
	return nil
}

// Any 注册所有 HTTP 方法的路由
func (group *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) error {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodOptions,
		http.MethodHead,
	}
	return group.Match(methods, relativePath, handlers...)
}

// combineHandlers 合并处理程序链
func (group *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	return append(group.handlers, handlers...) // 使用 append 合并处理程序
}

// calculateAbsolutePath 计算绝对路径
func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return osx.JoinPaths(group.basePath, relativePath) // 将基础路径与相对路径连接，返回绝对路径
}
