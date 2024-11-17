/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 22:55:15
 * @FilePath: \gosh\middlewares\recovery.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/kamalyes/gosh"
)

// RecoveryMiddleware 用于 recover 可能出现的 panic，并使用 zap 记录相关日志
func RecoveryMiddleware(stack bool) gosh.HandlerFunc {
	return func(c *gosh.Context) error {
		defer func() { handlePanic(c) }()
		c.Next()
		return nil
	}
}

func handlePanic(c *gosh.Context) error {
	if err := recover(); err != nil {
		httpRequest, _ := httputil.DumpRequest(c.Request, true)
		if isBrokenPipe(err) {
			c.Engine.Config.Zap.LogBrokenPipe(err, httpRequest)
			c.Abort()
			return err.(error)
		}
		c.Engine.Config.Zap.LogRecovery(err, httpRequest)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gosh.H{
			"message": "服务器内部错误，请稍后再试",
		})
	}
	return nil
}

func isBrokenPipe(err interface{}) bool {
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			errStr := strings.ToLower(se.Error())
			return strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer")
		}
	}
	return false
}
