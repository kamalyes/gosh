/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 16:06:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 11:50:52
 * @FilePath: \gosh\gosh_bench_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"net/http"
	"testing"
)

// 性能基准测试
func BenchmarkMainGoshOneRoute(b *testing.B) {
	// 每次基准测试前重新创建路由
	router := NewEngine()
	// 设置路由
	router.GET("/ping", func(c *Context) error {
		c.WriteString(http.StatusOK, "pong") // 返回简单字符串
		return nil
	})
	goshRunRequest(b, router, "GET", "/ping")
}

func BenchmarkGoshManyHandlers(B *testing.B) {
	router := NewEngine()
	router.Use(func(c *Context) error {
		// 处理逻辑
		return nil // 返回 nil 表示没有错误
	})
	router.Use(func(c *Context) error {
		// 处理逻辑
		return nil // 返回 nil 表示没有错误
	})
	router.GET("/ping", func(c *Context) error {
		c.WriteString(http.StatusOK, "pong") // 返回一个简单的字符串
		return nil                           // 返回 nil 表示没有错误
	})
	goshRunRequest(B, router, "GET", "/ping")
}

func BenchmarkGosh5Params(B *testing.B) {
	router := NewEngine()
	router.GET("/param/:param1/:param2/:param3/:param4/:param5", func(c *Context) error {
		c.WriteString(http.StatusOK, "Parameters received") // 返回一个简单的字符串
		return nil                                          // 返回 nil 表示没有错误
	})
	goshRunRequest(B, router, "GET", "/param/path/to/parameter/john/12345")
}

func BenchmarkGoshOneRouteJSON(B *testing.B) {
	router := NewEngine()
	data := struct {
		Status string `json:"status"`
	}{"ok"}
	router.GET("/json", func(c *Context) error {
		c.WriteJSONResponse(http.StatusOK, data)
		return nil // 返回 nil 表示没有错误
	})
	goshRunRequest(B, router, "GET", "/json")
}

func BenchmarkGoshOneRouteString(B *testing.B) {
	router := NewEngine()
	router.GET("/text", func(c *Context) error {
		c.WriteString(http.StatusOK, "this is a plain text")
		return nil // 返回 nil 表示没有错误
	})
	goshRunRequest(B, router, "GET", "/text")
}

func BenchmarkGoshManyRoutesLast(B *testing.B) {
	router := NewEngine()
	router.Any("/ping", func(c *Context) error {
		c.WriteString(http.StatusOK, "this is a plain text")
		return nil // 返回 nil 表示没有错误
	})
	goshRunRequest(B, router, "OPTIONS", "/ping")
}

func BenchmarkGosh404Many(B *testing.B) {
	router := NewEngine()
	router.GET("/", func(c *Context) error {
		c.WriteString(http.StatusOK, "Root") // 返回一个简单的字符串
		return nil                           // 返回 nil 表示没有错误
	})
	router.GET("/path/to/something", func(c *Context) error {
		c.WriteString(http.StatusOK, "Something") // 返回一个简单的字符串
		return nil                                // 返回 nil 表示没有错误
	})
	router.GET("/post/:id", func(c *Context) error {
		c.WriteString(http.StatusOK, "Post ID") // 返回一个简单的字符串
		return nil                              // 返回 nil 表示没有错误
	})
	router.GET("/view/:id", func(c *Context) error {
		c.WriteString(http.StatusOK, "View ID") // 返回一个简单的字符串
		return nil                              // 返回 nil 表示没有错误
	})
	router.GET("/favicon.ico", func(c *Context) error {
		c.WriteString(http.StatusOK, "Favicon") // 返回一个简单的字符串
		return nil                              // 返回 nil 表示没有错误
	})
	router.GET("/robots.txt", func(c *Context) error {
		c.WriteString(http.StatusOK, "Robots") // 返回一个简单的字符串
		return nil                             // 返回 nil 表示没有错误
	})
	router.GET("/delete/:id", func(c *Context) error {
		c.WriteString(http.StatusOK, "Delete ID") // 返回一个简单的字符串
		return nil                                // 返回 nil 表示没有错误
	})
	router.GET("/user/:id/:mode", func(c *Context) error {
		c.WriteString(http.StatusOK, "User  Mode") // 返回一个简单的字符串
		return nil                                 // 返回 nil 表示没有错误
	})

	goshRunRequest(B, router, "GET", "/viewfake")
}

type mockGoshWriter struct {
	headers http.Header
}

func newMockGoshWriter() *mockGoshWriter {
	return &mockGoshWriter{
		http.Header{},
	}
}

func (m *mockGoshWriter) Header() http.Header {
	return m.headers
}

func (m *mockGoshWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockGoshWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockGoshWriter) WriteHeader(int) {}

func goshRunRequest(B *testing.B, r *Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		B.Fatalf("Failed to create request: %v", err) // 使用 B.Fatalf 记录错误
	}
	w := newMockGoshWriter()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}
