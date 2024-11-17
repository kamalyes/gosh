/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 13:19:55
 * @FilePath: \gosh\zap.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	kmZap "github.com/kamalyes/go-config/pkg/zap"
	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/gosh/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogEncoderType 定义日志编码器类型
type LogEncoderType string

// 定义日志编码器类型常量，暴露给包外使用
const (
	LowercaseLevelEncoder      LogEncoderType = "LowercaseLevelEncoder"      // 小写级别编码器
	LowercaseColorLevelEncoder LogEncoderType = "LowercaseColorLevelEncoder" // 小写彩色级别编码器
	CapitalLevelEncoder        LogEncoderType = "CapitalLevelEncoder"        // 大写级别编码器
	CapitalColorLevelEncoder   LogEncoderType = "CapitalColorLevelEncoder"   // 大写彩色级别编码器
)

// String 方法实现，返回 LogEncoderType 的字符串表示
func (l LogEncoderType) String() string {
	return string(l)
}

// Logger 结构体封装了 zap.Logger 和一些自定义字段
type Logger struct {
	*zap.Logger             // 嵌入 zap.Logger 以便使用其方法
	kmZap         kmZap.Zap // 自定义的 zap 配置
	requestIDKey  string    // 请求 ID 键
	timeKey       string    // 时间键
	errorKey      string    // 错误键
	requestKey    string    // 请求键
	stacktraceKey string    // 堆栈跟踪键
}

// DefaultKmZipConfig 初始化默认的 KmZipConfig 配置
func DefaultKmZipConfig() kmZap.Zap {
	return kmZap.Zap{
		Director:      "./logs",                       // 默认日志目录
		EncodeLevel:   LowercaseLevelEncoder.String(), // 默认编码级别
		Format:        "json",                         // 默认格式为 JSON
		ShowLine:      true,                           // 默认显示行号
		StacktraceKey: "stacktrace",                   // 默认堆栈跟踪键
		MaxSize:       10,                             // 默认最大文件大小（MB）
		MaxBackups:    100,                            // 默认最大备份数量
		MaxAge:        30,                             // 默认最大保留天数
		Compress:      false,                          // 默认不压缩
		LogInConsole:  true,                           // 默认输出到控制台
	}
}

// NewLogger 初始化一个新的 Logger 实例
func NewLogger(ctx Context, kmZap kmZap.Zap) (*Logger, error) {
	// 检查日志目录是否存在
	if ok := osx.FileExists(kmZap.Director); !ok {
		fmt.Printf("创建目录：%v\n", kmZap.Director) // 输出创建目录信息
		if err := os.MkdirAll(kmZap.Director, os.ModePerm); err != nil {
			return nil, fmt.Errorf("创建目录失败: %v", err) // 创建目录失败时返回错误
		}
	}

	// 创建不同级别的核心
	cores := []zapcore.Core{
		getEncoderCore(ctx, zap.DebugLevel, kmZap), // 调试级别
		getEncoderCore(ctx, zap.InfoLevel, kmZap),  // 信息级别
		getEncoderCore(ctx, zap.WarnLevel, kmZap),  // 警告级别
		getEncoderCore(ctx, zap.ErrorLevel, kmZap), // 错误级别
	}

	// 使用 zapcore.NewTee 创建一个多核心的 logger
	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller())

	// 如果需要显示行号，则添加调用者信息
	if kmZap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	// 返回 Logger 实例
	return &Logger{
		Logger:        logger,
		kmZap:         kmZap,
		requestIDKey:  constants.LogRequestIDKey,
		timeKey:       constants.LogTimeKey,
		errorKey:      constants.LogErrorKey,
		requestKey:    constants.LogRequestKey,
		stacktraceKey: kmZap.StacktraceKey,
	}, nil
}

// getEncoderConfig 获取编码器配置
func getEncoderConfig(kmZap kmZap.Zap) zapcore.EncoderConfig {
	config := zapcore.EncoderConfig{
		MessageKey:    "message",                 // 消息键
		LevelKey:      "level",                   // 级别键
		TimeKey:       "time",                    // 时间键
		NameKey:       "logger",                  // 日志器名称键
		CallerKey:     "caller",                  // 调用者键
		StacktraceKey: kmZap.StacktraceKey,       // 堆栈跟踪键
		LineEnding:    zapcore.DefaultLineEnding, // 行结束符
		EncodeTime:    customTimeEncoder,         // 自定义时间编码器
		EncodeCaller:  zapcore.FullCallerEncoder, // 完整调用者编码器
	}

	// 根据配置选择编码级别
	switch kmZap.EncodeLevel {
	case LowercaseLevelEncoder.String():
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case LowercaseColorLevelEncoder.String():
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case CapitalLevelEncoder.String():
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case CapitalColorLevelEncoder.String():
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder // 默认使用小写级别编码器
	}

	return config
}

// getEncoder 获取编码器
func getEncoder(kmZap kmZap.Zap) zapcore.Encoder {
	// 根据格式返回相应的编码器
	if kmZap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(kmZap)) // JSON 编码器
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(kmZap)) // 控制台编码器
}

// getEncoderCore 获取编码器核心
func getEncoderCore(ctx Context, level zapcore.Level, kmZap kmZap.Zap) zapcore.Core {
	// 构建日志文件路径
	logFilePath := filepath.Join(kmZap.Director, fmt.Sprintf("%s-%s.log", ctx.Engine.Config.AppName, level.String()))
	writer := WriteSyncer(logFilePath, kmZap)                // 使用 WriteSyncer 创建写入器
	return zapcore.NewCore(getEncoder(kmZap), writer, level) // 创建并返回核心
}

// customTimeEncoder 自定义时间编码器
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339)) // 使用 RFC3339 格式化时间
}

// LogError 记录错误信息并返回 Logger 以支持链式调用
func (l *Logger) LogError(message string, err interface{}, httpRequest []byte, includeStack bool) *Logger {
	// 构建日志字段
	fields := []zap.Field{
		zap.String(l.requestIDKey, l.requestIDKey),    // 请求 ID
		zap.Time(l.timeKey, time.Now()),               // 当前时间
		zap.Any(l.errorKey, err),                      // 错误信息
		zap.String(l.requestKey, string(httpRequest)), // 请求内容
	}

	// 如果需要堆栈跟踪，则添加堆栈信息
	if includeStack {
		fields = append(fields, zap.Stack(l.stacktraceKey))
	}

	l.Error(message, fields...) // 记录错误日志
	return l                    // 返回 Logger 以支持链式调用
}

// LogBrokenPipe 记录断管错误
func (l *Logger) LogBrokenPipe(err interface{}, httpRequest []byte) *Logger {
	return l.LogError("broken pipe error", err, httpRequest, false) // 调用 LogError 记录断管错误
}

// LogRecovery 记录从 panic 恢复的信息
func (l *Logger) LogRecovery(err interface{}, httpRequest []byte) *Logger {
	return l.LogError("recovery from panic", err, httpRequest, true) // 调用 LogError 记录恢复信息
}

// WriteSyncer 利用 lumberjack 库做日志分割
func WriteSyncer(file string, kmZap kmZap.Zap) zapcore.WriteSyncer {
	// 日志文件的最大大小（以 MB 为单位）
	maxSize := kmZap.MaxSize
	if maxSize < 10 || maxSize > 500 {
		maxSize = 10 // 如果超出范围，则默认设置为 10MB
	}

	// 保留旧文件的最大个数
	maxBackups := kmZap.MaxBackups
	if maxBackups < 1 {
		maxBackups = 100 // 如果无效，则默认设置为 100
	}

	// 保留旧文件的最大天数
	maxAge := kmZap.MaxAge
	if maxAge < 1 {
		maxAge = 30 // 如果无效，则默认设置为 30 天
	}

	// 创建 lumberjack.Logger 实例
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,           // 日志文件路径
		MaxSize:    maxSize,        // 最大文件大小
		MaxBackups: maxBackups,     // 最大备份数量
		MaxAge:     maxAge,         // 最大保留天数
		Compress:   kmZap.Compress, // 是否压缩旧日志
	}

	// 如果需要在控制台输出日志
	if kmZap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger)) // 同时输出到控制台和文件
	}
	return zapcore.AddSync(lumberJackLogger) // 只输出到文件
}
