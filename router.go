/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 02:51:55
 * @FilePath: \gosh\router.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

// Router 路由器接口
type Router interface {
	Routes                                     // 嵌入 Routes 接口，表示 Router 也具有 Routes 的所有方法
	Group(string, ...HandlerFunc) *RouterGroup // 创建路由组的方法，返回一个 RouterGroup
}

// Routes 定义所有路由器接口
type Routes interface {
	Use(...HandlerFunc)                     // 注册中间件处理函数
	After(...HandlerFunc)                   // 在所有处理函数之后注册的处理函数
	Handle(string, string, ...HandlerFunc)  // 注册一个特定方法和路径的处理函数
	GET(string, ...HandlerFunc)             // 注册 GET 请求的处理函数
	POST(string, ...HandlerFunc)            // 注册 POST 请求的处理函数
	DELETE(string, ...HandlerFunc)          // 注册 DELETE 请求的处理函数
	PATCH(string, ...HandlerFunc)           // 注册 PATCH 请求的处理函数
	PUT(string, ...HandlerFunc)             // 注册 PUT 请求的处理函数
	OPTIONS(string, ...HandlerFunc)         // 注册 OPTIONS 请求的处理函数
	HEAD(string, ...HandlerFunc)            // 注册 HEAD 请求的处理函数
	Match([]string, string, ...HandlerFunc) // 根据请求方法数组和路径注册处理函数
	Any(string, ...HandlerFunc)             // 注册对任意请求方法的处理函数
	NoRoute(...HandlerFunc)                 // 注册没有匹配到路由时的处理函数
}
