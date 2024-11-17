/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 21:19:05
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 00:55:29
 * @FilePath: \gosh\errorsx\base.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package errorsx

import (
	"errors"
	"fmt"

	"github.com/kamalyes/gosh/constants"
)

// CustomError 是一个自定义错误类型，包含错误信息和错误类型
type CustomError struct {
	Err       error
	ErrorType ErrorType
}

// Error 实现 error 接口
func (e *CustomError) Error() string {
	return fmt.Sprintf("错误: %v, 类型: %d", e.Err, e.ErrorType)
}

// NewCustomError 创建一个新的 CustomError 实例
func NewCustomError(message string, errorType ErrorType) *CustomError {
	return &CustomError{
		Err:       errors.New(message),
		ErrorType: errorType,
	}
}

// 常用错误
var (
	ErrPathMustStartWithSlash    = NewCustomError(fmt.Sprintf("路径必须以%v开头", constants.PathSeparator), ErrorTypePublic)
	ErrMethodCannotBeEmpty       = NewCustomError("方法不能为空", ErrorTypePublic)
	ErrMustHaveAtLeastOneHandler = NewCustomError("必须有至少一个处理器", ErrorTypePublic)
	ErrWriteResponseFailed       = NewCustomError("写入响应时出错", ErrorTypePublic)
	ErrNotFound                  = NewCustomError("未找到请求的资源", ErrorTypePublic)
	ErrMethodNotAllowed          = NewCustomError("请求的方法不被允许", ErrorTypePublic)
	ErrInvalidRedirectCode       = NewCustomError("状态码必须在300到308之间", ErrorTypePublic)
	ErrNotMultipart              = NewCustomError("请求不是multipart格式", ErrorTypePublic)
	ErrAccessDenied              = NewCustomError("访问被拒绝", ErrorTypePublic)
	ErrFileNotFound              = NewCustomError("文件未找到", ErrorTypePublic)
	ErrInternalServerError       = NewCustomError("内部服务器错误", ErrorTypePublic)
	ErrDirectoryAccessForbidden  = NewCustomError("禁止访问目录", ErrorTypePublic)
)
