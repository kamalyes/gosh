/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 16:06:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 23:05:55
 * @FilePath: \gosh\reponse_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义常量
const (
	successMessage = "Success message"
	testData       = "test data"
)

// MockResponseWriter 是一个模拟的 ResponseWriter
type MockResponseWriter struct {
	StatusCode int
	Body       []byte
}

func (m *MockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *MockResponseWriter) Write(b []byte) (int, error) {
	m.Body = b
	return len(b), nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}

// setupContext 创建一个模拟的上下文
func setupContext() *Context {
	return &Context{
		ResponseWriter: &MockResponseWriter{}, // HTTP 响应写入器
	}
}

// assertResponseOption 检查 ResponseOption 的字段
func assertResponseOption(t *testing.T, resp *ResponseOption, expectedData interface{}, expectedCode SceneCode, expectedHttpCode StatusCode, expectedMessage string) {
	assert.Equal(t, expectedData, resp.Data)
	assert.Equal(t, expectedCode, resp.SceneCode)
	assert.Equal(t, expectedHttpCode, resp.HttpCode)
	assert.Equal(t, expectedMessage, resp.Message)
}

// TestNewResponseOption 测试 NewResponseOption 函数
func TestNewResponseOption(t *testing.T) {
	resp := NewResponseOption(testData, Fail, StatusOK, successMessage)
	assertResponseOption(t, resp, testData, Fail, StatusOK, successMessage)
}

// TestResponseOptionMerge 测试 merge 方法
func TestResponseOptionMerge(t *testing.T) {
	resp := &ResponseOption{
		SceneCode: 0,
		HttpCode:  0,
		Message:   "",
	}

	mergedResp := resp.Merge()
	assertResponseOption(t, mergedResp, nil, Success, Success, GetSceneCodeText(Success))
}

// TestSendJSONResponse 测试 SendJSONResponse 函数
func TestSendJSONResponse(t *testing.T) {
	ctx := setupContext()

	respOption := NewResponseOption(testData, Success, StatusOK, successMessage)
	err := SendJSONResponse(ctx, respOption)
	assert.NoError(t, err)
}

// TestGen400xResponse 测试 Gen400xResponse 函数
func TestGen400xResponse(t *testing.T) {
	ctx := setupContext()

	respOption := NewResponseOption(nil)
	Gen400xResponse(ctx, respOption)
	assertResponseOption(t, respOption, nil, BadRequest, StatusBadRequest, GetSceneCodeText(BadRequest))
}

// TestGen500xResponse 测试 Gen500xResponse 函数
func TestGen500xResponse(t *testing.T) {
	ctx := setupContext()

	respOption := NewResponseOption(nil)
	Gen500xResponse(ctx, respOption)
	assertResponseOption(t, respOption, nil, Fail, StatusInternalServerError, GetSceneCodeText(Fail))
}

// TestRemoveTopStruct 测试 removeTopStruct 函数
func TestRemoveTopStruct(t *testing.T) {
	fields := map[string]string{
		"User   .Name":  "Name is required",
		"User   .Email": "Email is required",
	}

	result := RemoveTopStruct(fields)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Name is required", result["Name"])
	assert.Equal(t, "Email is required", result["Email"])
}
