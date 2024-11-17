/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 16:06:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 17:10:22
 * @FilePath: \gosh\tests\reponse_test.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package tests

import (
	"net/http"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/kamalyes/gosh"
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
func setupContext() *gosh.Context {
	return &gosh.Context{
		ResponseWriter: &MockResponseWriter{}, // HTTP 响应写入器
	}
}

// assertResponseOption 检查 ResponseOption 的字段
func assertResponseOption(t *testing.T, resp *gosh.ResponseOption, expectedData interface{}, expectedCode gosh.SceneCode, expectedHttpCode gosh.StatusCode, expectedMessage string) {
	assert.Equal(t, expectedData, resp.Data)
	assert.Equal(t, expectedCode, resp.SceneCode)
	assert.Equal(t, expectedHttpCode, resp.HttpCode)
	assert.Equal(t, expectedMessage, resp.Message)
}

// TestNewResponseOption 测试 gosh.NewResponseOption 函数
func TestNewResponseOption(t *testing.T) {
	resp := gosh.NewResponseOption(testData, gosh.Fail, gosh.StatusOK, successMessage)
	assertResponseOption(t, resp, testData, gosh.Fail, gosh.StatusOK, successMessage)
}

// TestResponseOptionMerge 测试 merge 方法
func TestResponseOptionMerge(t *testing.T) {
	resp := &gosh.ResponseOption{
		SceneCode: 0,
		HttpCode:  0,
		Message:   "",
	}

	mergedResp := resp.Merge()
	assertResponseOption(t, mergedResp, nil, gosh.Success, gosh.Success, gosh.GetSceneCodeText(gosh.Success))
}

// TestSendJSONResponse 测试 SendJSONResponse 函数
func TestSendJSONResponse(t *testing.T) {
	ctx := setupContext()

	respOption := gosh.NewResponseOption(testData, gosh.Success, gosh.StatusOK, successMessage)
	err := gosh.SendJSONResponse(ctx, respOption)
	assert.NoError(t, err)
}

// TestGen400xResponse 测试 Gen400xResponse 函数
func TestGen400xResponse(t *testing.T) {
	ctx := setupContext()

	respOption := gosh.NewResponseOption(nil)
	gosh.Gen400xResponse(ctx, respOption)
	assertResponseOption(t, respOption, nil, gosh.BadRequest, gosh.StatusBadRequest, gosh.GetSceneCodeText(gosh.BadRequest))
}

// TestGen500xResponse 测试 Gen500xResponse 函数
func TestGen500xResponse(t *testing.T) {
	ctx := setupContext()

	respOption := gosh.NewResponseOption(nil)
	gosh.Gen500xResponse(ctx, respOption)
	assertResponseOption(t, respOption, nil, gosh.Fail, gosh.StatusInternalServerError, gosh.GetSceneCodeText(gosh.Fail))
}

// TestValidatorError 测试 ValidatorError 函数
func TestValidatorError(t *testing.T) {
	ctx := setupContext()

	// 假设这里是一个验证错误
	validationErr := validator.ValidationErrors{}
	gosh.ValidatorError(ctx, validationErr)

	// 根据您的实现，您可以添加适当的断言
}

// TestRemoveTopStruct 测试 removeTopStruct 函数
func TestRemoveTopStruct(t *testing.T) {
	fields := map[string]string{
		"User   .Name":  "Name is required",
		"User   .Email": "Email is required",
	}

	result := gosh.RemoveTopStruct(fields)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Name is required", result["Name"])
	assert.Equal(t, "Email is required", result["Email"])
}
