/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 17:51:15
 * @FilePath: \gosh\banner_test.go
 * @Description: 测试 Banner 功能
 */

package gosh

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"

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
	config := NewBannerConfig().SetBanner("Test Banner").SetTitle("Test Title").SetSubtitle("Test Subtitle")

	// 打印 Banner
	config.Print()

	// 关闭写入端，准备读取输出
	w.Close()
	var outputBuf bytes.Buffer
	io.Copy(&outputBuf, r) // 从管道读取输出到 outputBuf

	// 恢复标准输出
	os.Stdout = oldStdout

	// 使用 assert 检查输出是否符合预期
	expectedOutput := fmt.Sprintf("%s\n%s\n%s\n\n", config.GetBanner(), config.GetTitle(), config.GetSubtitle())
	assert.Equal(t, expectedOutput, outputBuf.String(), "Output should match expected format")
}

// 测试默认值
func TestDefaultBannerConfig(t *testing.T) {
	config := NewBannerConfig()

	assert.Equal(t, constants.DefaultBanner, config.GetBanner(), "Default banner should match")
	assert.Equal(t, constants.DefaultTitle, config.GetTitle(), "Default title should match")
	assert.Equal(t, constants.DefaultSubtitle, config.GetSubtitle(), "Default subtitle should match")
}

// 测试设置 Banner
func TestSetBanner(t *testing.T) {
	config := NewBannerConfig()
	config.SetBanner("New Banner")

	assert.Equal(t, "New Banner", config.GetBanner(), "Banner should be updated")
}

// 测试设置 Title
func TestSetTitle(t *testing.T) {
	config := NewBannerConfig()
	config.SetTitle("New Title")

	assert.Equal(t, "New Title", config.GetTitle(), "Title should be updated")
}

// 测试设置 Subtitle
func TestSetSubtitle(t *testing.T) {
	config := NewBannerConfig()
	config.SetSubtitle("New Subtitle")

	assert.Equal(t, "New Subtitle", config.GetSubtitle(), "Subtitle should be updated")
}

// 测试并发访问
func TestConcurrentSetAndPrint(t *testing.T) {
	config := NewBannerConfig()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		config.SetBanner("Concurrent Banner")
	}()

	go func() {
		defer wg.Done()
		config.Print()
	}()

	wg.Wait()
}
