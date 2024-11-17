/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 01:37:25
 * @FilePath: \gosh\tests\static_test.go
 * @Description: 测试 StaticFile 功能
 */
package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/kamalyes/gosh" // 确保导入您的框架
	"github.com/stretchr/testify/assert"
)

// createStaticTempFile 创建一个临时文件并写入内容
func createStaticTempFile(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "testfile-*.txt")
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

// TestStatic 测试 Static 方法
func TestStatic(t *testing.T) {
	// 创建一个临时目录
	tempDir := t.TempDir()

	// 创建两个临时文件
	file1Content := "Hello from file 1"
	file1Path, err := createStaticTempFile(file1Content)
	assert.NoError(t, err, "Failed to create temp file 1")
	defer os.Remove(file1Path) // 确保测试结束后删除文件

	file2Content := "Hello from file 2"
	file2Path, err := createStaticTempFile(file2Content)
	assert.NoError(t, err, "Failed to create temp file 2")
	defer os.Remove(file2Path) // 确保测试结束后删除文件

	// 将文件移动到临时目录
	newFile1Path := filepath.Join(tempDir, "file1.txt")
	os.Rename(file1Path, newFile1Path)
	fmt.Println("File 1 moved to:", newFile1Path) // 打印文件路径

	newFile2Path := filepath.Join(tempDir, "file2.txt")
	os.Rename(file2Path, newFile2Path)
	fmt.Println("File 2 moved to:", newFile2Path) // 打印文件路径

	// 创建一个新的 RouterGroup
	engine := gosh.NewEngine()
	group := engine.Group("/static")

	// 注册静态文件路由
	group.Static("/files", tempDir, true)
	fmt.Println("Static files served from:", tempDir) // 打印静态文件服务路径

	// 创建一个测试请求
	req1, err := http.NewRequest(http.MethodGet, "/static/files/file1.txt", nil)
	assert.NoError(t, err, "Failed to create request for file1")
	fmt.Println("Requesting:", req1.URL.Path) // 打印请求路径

	// 创建一个响应记录器
	rr1 := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(rr1, req1)

	// 检查响应状态码和内容
	if rr1.Code != http.StatusOK {
		fmt.Println("Error serving file 1:", rr1.Body.String()) // 打印错误信息
	}
	assert.Equal(t, http.StatusOK, rr1.Code, "Handler returned wrong status code for file1")
	assert.Equal(t, file1Content, rr1.Body.String(), "Handler returned unexpected body for file1")

	// 测试另一个文件
	req2, err := http.NewRequest(http.MethodGet, "/static/files/file2.txt", nil)
	assert.NoError(t, err, "Failed to create request for file2")
	fmt.Println("Requesting:", req2.URL.Path) // 打印请求路径

	rr2 := httptest.NewRecorder()
	engine.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		fmt.Println("Error serving file 2:", rr2.Body.String()) // 打印错误信息
	}
	assert.Equal(t, http.StatusOK, rr2.Code, "Handler returned wrong status code for file2")
	assert.Equal(t, file2Content, rr2.Body.String(), "Handler returned unexpected body for file2")

	// 测试不存在的文件
	req3, err := http.NewRequest(http.MethodGet, "/static/files/nonexistent.txt", nil)
	assert.NoError(t, err, "Failed to create request for nonexistent file")

	rr3 := httptest.NewRecorder()
	engine.ServeHTTP(rr3, req3)

	assert.Equal(t, http.StatusNotFound, rr3.Code, "Handler should return 404 for nonexistent file")
}
