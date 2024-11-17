/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 22:36:15
 * @FilePath: \gosh\middlewares\trace.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"context"

	"github.com/kamalyes/go-toolbox/pkg/osx"
	"github.com/kamalyes/gosh"
	"github.com/kamalyes/gosh/constants"
)

// NewTraceIDContext 将追踪ID存储到上下文中
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, constants.TraceRequestKey, traceID)
}

// GetTraceID 从上下文中获取追踪ID
func GetTraceID(ctx context.Context) (string, bool) {
	if v := ctx.Value(constants.TraceRequestKey); v != nil {
		if id, ok := v.(string); ok && id != "" {
			return id, true
		}
	}
	return "", false
}

// TraceMiddleware 追踪中间件处理函数
func TraceMiddleware(skippers ...SkipperFunc) gosh.HandlerFunc {
	return func(c *gosh.Context) error {
		// 检查是否跳过中间件
		if SkipHandler(c, skippers...) {
			c.Next()
			return nil
		}

		// 从请求头中获取追踪ID，若为空则生成一个新的追踪ID
		traceID := c.Header(constants.TraceRequestKey)
		if traceID == "" {
			traceID = osx.HashUnixMicroCipherText()
		}

		// 将追踪ID存储到上下文中，并设置响应头中的追踪ID
		ctx := NewTraceIDContext(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)
		c.ResponseWriter.Header().Set(constants.TraceRequestKey, traceID)
		c.Next()
		return nil
	}
}
