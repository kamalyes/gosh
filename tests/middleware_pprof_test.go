/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 18:56:57
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 22:55:55
 * @FilePath: \gosh\tests\middleware_pprof_test.go
 * @Description: 对 pprof 包进行单元测试
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/kamalyes/gosh"
	"github.com/kamalyes/gosh/middlewares"
	"github.com/stretchr/testify/assert"
)

func TestNewSystemInfo(t *testing.T) {
	// 创建 SystemInfo 实例
	sysInfo := middlewares.NewSystemInfo()

	// 使用 assert 来验证各字段
	assert.NotEmpty(t, sysInfo.ServerName, "ServerName should not be empty")
	assert.NotEmpty(t, sysInfo.Runtime, "Runtime should not be empty")
	assert.NotEmpty(t, sysInfo.GoroutineNum, "GoroutineNum should not be empty")
	assert.NotEmpty(t, sysInfo.CPUNum, "CPUNum should not be empty")
	assert.NotEmpty(t, sysInfo.UsedMem, "UsedMem should not be empty")
	assert.NotEmpty(t, sysInfo.HeapInuse, "HeapInuse should not be empty")
	assert.NotEmpty(t, sysInfo.TotalMem, "TotalMem should not be empty")
	assert.NotEmpty(t, sysInfo.SysMem, "SysMem should not be empty")
	assert.NotEmpty(t, sysInfo.Lookups, "Lookups should not be empty")
	assert.NotEmpty(t, sysInfo.Mallocs, "Mallocs should not be empty")
	assert.NotEmpty(t, sysInfo.Frees, "Frees should not be empty")
	assert.NotEmpty(t, sysInfo.LastGCTime, "LastGCTime should not be empty")
	assert.NotEmpty(t, sysInfo.NextGC, "NextGC should not be empty")
	assert.NotEmpty(t, sysInfo.PauseTotalNs, "PauseTotalNs should not be empty")
	assert.NotEmpty(t, sysInfo.PauseNs, "PauseNs should not be empty")
}

func TestHandler(t *testing.T) {
	// 创建一个 HTTP 测试请求和响应
	req := httptest.NewRequest("GET", "/debug/pprof/sysinfo", nil)
	w := httptest.NewRecorder()

	// 创建一个新的上下文
	ctx := &gosh.Context{
		Request:        req,
		ResponseWriter: w,
	}

	// 调用 Handler
	err := middlewares.PprofHandler(ctx)
	assert.NoError(t, err, "Handler should not return an error")

	// 检查响应状态码
	assert.Equal(t, 200, w.Code, "Expected status code 200")

	// 检查响应内容是否包含关键字（例如服务器名称）
	assert.Contains(t, w.Body.String(), "服务器", "Response should contain '服务器'")
	assert.Contains(t, w.Body.String(), "运行时间", "Response should contain '运行时间'")
	assert.Contains(t, w.Body.String(), "goroutine数量", "Response should contain 'goroutine数量'")
	assert.Contains(t, w.Body.String(), "CPU核数", "Response should contain 'CPU核数'")
}
