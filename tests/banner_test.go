/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 13:55:55
 * @FilePath: \gosh\tests\banner_test.go
 * @Description: 测试 Banner 功能
 */

package tests

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/kamalyes/gosh"
	"github.com/kamalyes/gosh/constants"
	"github.com/stretchr/testify/assert"
)

// 测试 PrintBanner
func TestPrintBanner(t *testing.T) {
	// 保存原始的 os.Stdout
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }() // 恢复原始的 os.Stdout

	// 将 os.Stdout 重定向到 buf
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 创建一个测试配置
	config := gosh.BannerConfig{
		Banner:   "Test Banner",
		Title:    "Test Title",
		Subtitle: "Test Subtitle",
	}

	// 打印 Banner
	gosh.PrintBanner(config)

	// 关闭写入端，准备读取输出
	w.Close()
	var outputBuf bytes.Buffer
	io.Copy(&outputBuf, r) // 从管道读取输出到 outputBuf

	// 恢复标准输出
	os.Stdout = oldStdout

	// 使用 assert 检查输出是否符合预期
	expectedOutput := fmt.Sprintf("%s\n   %s\n   %s\n\n", config.Banner, config.Title, config.Subtitle)
	assert.Equal(t, expectedOutput, outputBuf.String(), "Output should match expected format")
}

// 测试 GetDefaultBanner
func TestGetDefaultBanner(t *testing.T) {
	// 获取默认 Banner 配置
	defaultBanner := gosh.DefaultBannerConfig()

	// 使用 assert 验证默认 Banner 配置是否符合预期
	assert.Equal(t, constants.DefaultBanner, defaultBanner.Banner, "Banner should match default value")
	assert.Equal(t, constants.DefaultTitle, defaultBanner.Title, "Title should match default value")
	assert.Equal(t, constants.DefaultSubtitle, defaultBanner.Subtitle, "Subtitle should match default value")
}
