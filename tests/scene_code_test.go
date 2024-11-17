/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 17:16:34
 * @FilePath: \gosh\tests\scene_code_test.go
 * @Description: 测试 SceneCode 功能
 */

package tests

import (
	"testing"

	"github.com/kamalyes/gosh"
)

// 测试 SetSceneCode 和 GetSceneCodeText
func TestSetAndGetSceneCode(t *testing.T) {
	// 定义一个自定义状态码
	const customCode gosh.SceneCode = 2001
	customMsg := "Custom Success"

	// 设置自定义状态码
	gosh.SetSceneCode(customCode, customMsg)

	// 获取自定义状态码对应的消息
	got := gosh.GetSceneCodeText(customCode)
	if got != customMsg {
		t.Errorf("expected %s, got %s", customMsg, got)
	}
}

// 测试 InitSceneCodes
func TestInitSceneCodes(t *testing.T) {
	// 定义自定义状态码和消息
	customCodes := map[gosh.SceneCode]string{
		200:  "Custom Success",
		2002: "Custom Error",
		2003: "Another Custom Error",
	}

	// 初始化自定义状态码
	gosh.InitSceneCodes(customCodes)

	// 验证自定义状态码的消息
	for code, expectedMsg := range customCodes {
		got := gosh.GetSceneCodeText(code)
		if got != expectedMsg {
			t.Errorf("expected %s for code %d, got %s", expectedMsg, code, got)
		}
	}
}

// 测试获取未设置的状态码消息
func TestGetUnknownSceneCode(t *testing.T) {
	unknownCode := gosh.SceneCode(9999) // 一个未定义的状态码

	// 获取未知状态码的消息
	got := gosh.GetSceneCodeText(unknownCode)
	if got != "" {
		t.Errorf("expected empty string for unknown code, got %s", got)
	}
}
