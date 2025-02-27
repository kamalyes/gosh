/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 18:56:57
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 23:55:56
 * @FilePath: \gosh\example_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// 错误回调处理器
func errorHandler(ctx *Context) {
	// 记录错误日志
	log.Println("记录错误日志", ctx.Status, ctx.Error)
	// 输出错误信息到客户端
	ctx.ResponseWriter.WriteHeader(ctx.Status)
	if ctx.Error != nil {
		_, _ = ctx.ResponseWriter.Write([]byte(ctx.Error.Error()))
	}
}

// 后置回调处理器
func afterHandler(ctx *Context) {
	log.Println("执行了后置处理器", ctx.IsAborted())
}

// 测试回应
func TestEcho(t *testing.T) {
	app := NewEngine()
	app.GET("/", func(ctx *Context) error {
		t.Log("Hello Tsing")
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

func TestStatusCode(t *testing.T) {
	app := NewEngine()
	app.GET("/", func(ctx *Context) error {
		return ctx.WriteNoContent()
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	resp := httptest.NewRecorder()
	app.ServeHTTP(resp, req)
	t.Log(resp.Code)
}

// 测试处理器
func TestHandlers(t *testing.T) {
	app := NewEngine(Config{
		Recovery:     true,
		AfterHandler: afterHandler,
	})
	app.Use(func(ctx *Context) error {
		t.Log("1 执行了全局中间件")
		return nil
	})
	group := app.Group("/group", func(ctx *Context) error {
		t.Log("2 执行了 /group")
		return nil
	})
	group.Use(func(ctx *Context) error {
		t.Log("3 执行了路由组 /group 中间件")
		return nil
	})
	group.GET("/object", func(ctx *Context) error {
		t.Log("4 执行了 /group/object")
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "/group2/object", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试 PathValue
func TestPathValue(t *testing.T) {
	app := NewEngine()
	app.GET("/:path/:file", func(ctx *Context) error {
		t.Log("path=", ctx.PathValue("path"))
		t.Log("file=", ctx.PathValue("file"))
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "/haha/123", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试Context传值
func TestContextValue(t *testing.T) {
	app := NewEngine()
	app.GET("/", func(ctx *Context) error {
		// 在ctx中写入参数
		ctx.SetContextValue("hello", "tsing")
		return nil
	}, func(ctx *Context) error {
		t.Log("hello=", ctx.GetContextValue("hello"))
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试中止处理器链
func TestAbort(t *testing.T) {
	app := NewEngine()
	group := app.Group("/group")
	group.GET("/object", func(ctx *Context) error {
		t.Log("ok")
		ctx.Abort()
		return nil
	}, func(ctx *Context) error {
		t.Error("中止失败")
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "/group/object", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试QueryValue
func TestQueryParams(t *testing.T) {
	app := NewEngine()
	app.GET("/", func(ctx *Context) error {
		t.Log("id=", ctx.QueryValue("id"))
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "/?id=123", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试FormValue
func TestFormValue(t *testing.T) {
	app := NewEngine()
	app.POST("/", func(ctx *Context) error {
		t.Log("test=", ctx.FormValue("test"))
		return nil
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "POST", "/", strings.NewReader("test=ok"))
	if err != nil {
		t.Error(err)
		return
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试 404 错误
func TestNotFoundError(t *testing.T) {
	app := NewEngine(Config{
		ErrorHandler: errorHandler,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "/404", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试 405 事件
func TestMethodNotAllowedError(t *testing.T) {
	app := NewEngine(Config{
		HandleMethodNotAllowed: true,
		ErrorHandler:           errorHandler,
	})
	app.POST("/", func(ctx *Context) error {
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 测试CORS
func TestCORS(t *testing.T) {
	app := NewEngine(Config{
		ErrorHandler:           errorHandler, // 通过错误处理器来实现自动响应OPTIONS请求
		HandleMethodNotAllowed: true,         // 错误处理器中需要判断 405 状态码
	})
	app.GET("/", func(ctx *Context) error {
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := http.NewRequestWithContext(ctx, "OPTIONS", "/", nil)
	if err != nil {
		t.Error(err)
		return
	}
	app.ServeHTTP(httptest.NewRecorder(), r)
}

// 前置处理器的测试变量
var logOutput string

// 前置处理器示例
func beforeHandler(ctx *Context) {
	logOutput += "前置处理器: 处理请求 " + ctx.Request.Method + " " + ctx.Request.URL.Path + "\n"
	log.Println(logOutput)
}

// 测试前置处理器
func TestBeforeHandler(t *testing.T) {
	// 创建引擎实例并设置前置处理器
	engine := NewEngine(Config{
		BeforeHandler: beforeHandler,
	})

	// 设置路由
	engine.GET("/", func(ctx *Context) error {
		ctx.ResponseWriter.Write([]byte("Hello, World!"))
		return nil
	})

	// 创建请求
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}

	// 创建响应记录器
	rr := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(rr, req)

	// 检查响应状态码
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("期望状态码 %v，但得到 %v", http.StatusOK, status)
	}

	// 检查响应内容
	expected := "Hello, World!"
	if rr.Body.String() != expected {
		t.Errorf("期望响应内容 %q，但得到 %q", expected, rr.Body.String())
	}

	// 检查前置处理器的输出
	expectedLog := "前置处理器: 处理请求 GET /\n"
	if logOutput != expectedLog {
		t.Errorf("期望前置处理器输出 %q，但得到 %q", expectedLog, logOutput)
	}
}
