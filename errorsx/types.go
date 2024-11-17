/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 21:19:05
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 02:42:25
 * @FilePath: \gosh\errorsx\types.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package errorsx

// ErrorType 是一个无符号 64 位错误代码，按照 gosh 规范定义。
type ErrorType uint64

const (
	// ErrorTypeBind 用于表示 Context.Bind() 失败时的错误类型。
	ErrorTypeBind ErrorType = 1 << 63
	// ErrorTypeRender 用于表示 Context.Render() 失败时的错误类型。
	ErrorTypeRender ErrorType = 1 << 62
	// ErrorTypePrivate 表示一个私有错误。
	ErrorTypePrivate ErrorType = 1 << 0
	// ErrorTypePublic 表示一个公共错误。
	ErrorTypePublic ErrorType = 1 << 1
	// ErrorTypeAny 表示任何其他错误。
	ErrorTypeAny ErrorType = 1<<64 - 1
	// ErrorTypeNu 表示另一种错误类型。
	ErrorTypeNu = 2
)
