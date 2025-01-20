/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 18:16:30
 * @FilePath: \gosh\router_group_test.go
 * @Description: 测试 RouterGroup 功能
 */

package gosh

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 RouterGroup 的基本功能
func TestRouterGroup(t *testing.T) {
	engine := NewEngine()         // 假设您有一个创建 Engine 的方法
	group := engine.Group("/api") // 创建一个路由组

	// 注册一个简单的处理程序
	group.GET("/hello", func(c *Context) error {
		c.WriteString(http.StatusOK, "Hello, World!")
		return nil
	})

	// 创建一个测试请求
	req, err := http.NewRequest(http.MethodGet, "/api/hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 创建一个响应记录器
	recorder := httptest.NewRecorder()

	// 执行请求
	engine.ServeHTTP(recorder, req)

	// 断言响应状态码和内容
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "Hello, World!", recorder.Body.String())
}

// 测试中间件的使用
func TestRouterGroupMiddleware(t *testing.T) {
	engine := NewEngine()
	group := engine.Group("/api")

	// 注册中间件
	group.Use(func(c *Context) error {
		c.SetContextValue("key", "value") // 设置上下文中的值
		c.Next()                          // 调用下一层处理程序
		return nil
	})

	// 注册处理程序
	group.GET("/test", func(c *Context) error {
		value := c.GetContextValue("key")            // 获取上下文中的值
		c.WriteString(http.StatusOK, value.(string)) // 返回该值
		return nil
	})

	// 创建一个测试请求
	req, err := http.NewRequest(http.MethodGet, "/api/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	// 断言响应状态码和内容
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "value", recorder.Body.String())
}

func TestRouterGroup_ConflictHandling(t *testing.T) {
	engine := NewEngine()
	group := &RouterGroup{Engine: engine}

	// 定义一个简单的处理函数
	handler1 := func(ctx *Context) error {
		// 处理逻辑
		ctx.ResponseWriter.Write([]byte("Hello, Handler1"))
		return nil // 或者返回一个错误
	}

	handler2 := func(ctx *Context) error {
		// 处理逻辑
		ctx.ResponseWriter.Write([]byte("Hello, Handler2"))
		return nil // 或者返回一个错误
	}

	// 注册不同方法的路由
	group.GET("/test", handler1)
	group.POST("/test", handler2)
	group.PUT("/test", handler2)
	group.PATCH("/test", handler2)

	// 尝试注册相同路径的不同方法并捕获 panic
	assert.Panics(t, func() {
		group.GET("/test", handler2) // 这应该导致 panic
	}, "Expected panic when registering the same route")
}
