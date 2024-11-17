/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 18:27:17
 * @FilePath: \gosh\tests\static_file_test.go
 * @Description: 测试 StaticFile 功能
 */

package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/kamalyes/gosh"
	"github.com/stretchr/testify/assert"
)

// createStaticFileTempFile 创建一个临时文件并写入内容
func createStaticFileTempFile(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		return "", err
	}
	_, err = tempFile.Write([]byte(content))
	if err != nil {
		return "", err
	}
	tempFile.Close() // 关闭文件以确保可以读取
	return tempFile.Name(), nil
}

// testStaticFileHelper 是一个辅助函数，用于测试静态文件路由
func testStaticFileHelper(t *testing.T, route string, content string, isFS bool) {
	// 创建一个临时文件
	tempFileName, err := createStaticFileTempFile(content)
	assert.NoError(t, err, "Failed to create temp file")
	defer os.Remove(tempFileName) // 测试结束后删除临时文件

	// 获取临时文件的目录
	tempDir := filepath.Dir(tempFileName)

	// 创建一个新的 RouterGroup
	engine := gosh.NewEngine()       // 创建 Engine
	group := engine.Group("/static") // 创建一个路由组

	// 根据 isFS 参数选择注册静态文件的方法
	if isFS {
		// 使用自定义文件系统注册静态文件路由
		group.StaticFileFS(route, filepath.Base(tempFileName), http.Dir(tempDir)) // 使用临时文件的基础名称
	} else {
		// 注册静态文件路由
		group.StaticFile(route, tempFileName)
	}

	// 创建一个测试请求
	req, err := http.NewRequest(http.MethodGet, "/static"+route, nil)
	assert.NoError(t, err, "Failed to create request")

	// 创建一个响应记录器
	rr := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(rr, req)

	// 检查响应状态码
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	// 检查响应内容
	expected := content
	assert.Equal(t, expected, rr.Body.String(), "Handler returned unexpected body")
}

// TestStaticFile 测试 StaticFile 方法
func TestStaticFile(t *testing.T) {
	testStaticFileHelper(t, "/testfile", "Hello, this is a test file.", false)
}

// TestStaticFileFS 测试 StaticFileFS 方法
func TestStaticFileFS(t *testing.T) {
	testStaticFileHelper(t, "/testfilefs", "Hello from FS test file.", true)
}
