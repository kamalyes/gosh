/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 17:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 13:25:06
 * @FilePath: \gosh\tests\zap_bench_test.go
 * @Description: 测试 Zap 性能
 */
package tests

import (
	"errors"
	"testing"
)

// 基准测试 LogError 方法
func BenchmarkLogError(b *testing.B) {
	logger, _ := initLogger(&testing.T{}) // 使用辅助函数初始化 Logger
	testErr := errors.New("测试错误")         // 创建一个测试错误
	httpRequest := []byte("测试请求")         // 创建一个测试请求内容

	for i := 0; i < b.N; i++ {
		logger.LogError("测试错误信息", testErr, httpRequest, true)
	}
}

// 基准测试 LogBrokenPipe 方法
func BenchmarkLogBrokenPipe(b *testing.B) {
	logger, _ := initLogger(&testing.T{}) // 使用辅助函数初始化 Logger
	testErr := errors.New("断管错误")         // 创建一个测试断管错误
	httpRequest := []byte("测试请求")         // 创建一个测试请求内容

	for i := 0; i < b.N; i++ {
		logger.LogBrokenPipe(testErr, httpRequest)
	}
}

// 基准测试 LogRecovery 方法
func BenchmarkLogRecovery(b *testing.B) {
	logger, _ := initLogger(&testing.T{}) // 使用辅助函数初始化 Logger
	testErr := errors.New("恢复错误")         // 创建一个测试恢复错误
	httpRequest := []byte("测试请求")         // 创建一个测试请求内容

	for i := 0; i < b.N; i++ {
		logger.LogRecovery(testErr, httpRequest)
	}
}
