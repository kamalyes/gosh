/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 13:26:42
 * @FilePath: \gosh\tests\zap_test.go
 * @Description: 测试 Zap 功能
 */
package tests

import (
	"errors"
	"testing"

	"github.com/kamalyes/go-config/pkg/zap"
	"github.com/kamalyes/gosh"
	"github.com/stretchr/testify/assert"
)

// 初始化 Logger 的辅助函数
func initLogger(t *testing.T) (*gosh.Logger, zap.Zap) {
	kmZap := gosh.DefaultKmZipConfig() // 创建一个默认的 kmZap 配置
	ctx := gosh.Context{
		Engine: &gosh.Engine{
			Config: gosh.Config{
				AppName: "testate",
			},
		},
	}
	logger, err := gosh.NewLogger(ctx, kmZap) // 初始化一个新的 Logger 实例
	assert.NoError(t, err)                    // 断言没有错误发生
	assert.NotNil(t, logger)                  // 断言 logger 不为 nil
	return logger, kmZap
}

// TestNewLogger 测试 NewLogger 方法
func TestNewLogger(t *testing.T) {
	_, _ = initLogger(t) // 使用辅助函数初始化 Logger
}

// TestLogError 测试 LogError 方法
func TestLogError(t *testing.T) {
	logger, _ := initLogger(t) // 使用辅助函数初始化 Logger

	testErr := errors.New("测试错误") // 创建一个测试错误
	httpRequest := []byte("测试请求") // 创建一个测试请求内容
	logger.LogError("测试错误信息", testErr, httpRequest, true)

	// 在这里可以添加断言来验证日志内容是否正确
}

// TestLogBrokenPipe 测试 LogBrokenPipe 方法
func TestLogBrokenPipe(t *testing.T) {
	logger, _ := initLogger(t) // 使用辅助函数初始化 Logger

	testErr := errors.New("断管错误") // 创建一个测试断管错误
	httpRequest := []byte("测试请求") // 创建一个测试请求内容
	logger.LogBrokenPipe(testErr, httpRequest)

	// 在这里可以添加断言来验证日志内容是否正确
}

// TestLogRecovery 测试 LogRecovery 方法
func TestLogRecovery(t *testing.T) {
	logger, _ := initLogger(t) // 使用辅助函数初始化 Logger

	testErr := errors.New("恢复错误") // 创建一个测试恢复错误
	httpRequest := []byte("测试请求") // 创建一个测试请求内容
	logger.LogRecovery(testErr, httpRequest)

	// 在这里可以添加断言来验证日志内容是否正确
}

// TestWriteSyncer 测试 WriteSyncer 方法
func TestWriteSyncer(t *testing.T) {
	kmZap := gosh.DefaultKmZipConfig()
	filePath := "test.log"                           // 测试日志文件路径
	writeSyncer := gosh.WriteSyncer(filePath, kmZap) // 创建 WriteSyncer 实例
	assert.NotNil(t, writeSyncer)                    // 断言 WriteSyncer 不为 nil
}
