/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 17:17:08
 * @FilePath: \gosh\status_code_test.go
 * @Description: 测试 StatusCode 功能
 */

package gosh

import (
	"testing"
)

// 测试 SetStatusCode 和 GetStatusCodeText
func TestSetAndGetStatusCode(t *testing.T) {
	// 定义一个自定义状态码
	const customCode StatusCode = 2001
	customMsg := "Custom Success"

	// 设置自定义状态码
	SetStatusCode(customCode, customMsg)

	// 获取自定义状态码对应的消息
	got := GetStatusCodeText(customCode)
	if got != customMsg {
		t.Errorf("expected %s, got %s", customMsg, got)
	}
}

// 测试 InitStatusCodes
func TestInitStatusCodes(t *testing.T) {
	// 定义自定义状态码和消息
	customCodes := map[StatusCode]string{
		200:  "Custom Success",
		2002: "Custom Error",
		2003: "Another Custom Error",
	}

	// 初始化自定义状态码
	InitStatusCodes(customCodes)

	// 验证自定义状态码的消息
	for code, expectedMsg := range customCodes {
		got := GetStatusCodeText(code)
		if got != expectedMsg {
			t.Errorf("expected %s for code %d, got %s", expectedMsg, code, got)
		}
	}
}

// 测试获取未设置的状态码消息
func TestGetUnknownStatusCode(t *testing.T) {
	unknownCode := StatusCode(9999) // 一个未定义的状态码

	// 获取未知状态码的消息
	got := GetStatusCodeText(unknownCode)
	if got != "" {
		t.Errorf("expected empty string for unknown code, got %s", got)
	}
}
