/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 23:05:55
 * @FilePath: \gosh\middlewares\body.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"sync"

	"github.com/kamalyes/gosh"
)

var (
	maxMemory   int64 = 64 << 20 // 64 MB
	memoryMutex sync.Mutex
	reqBodyKey  string = "/req-body"
)

// GetMaxMemory 获取最大内存
func GetMaxMemory() int64 {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()
	return maxMemory
}

// SetMaxMemory 设置最大内存
func SetMaxMemory(memory int64) {
	memoryMutex.Lock()
	defer memoryMutex.Unlock()
	maxMemory = memory
}

func CopyBodyMiddleware(skippers ...SkipperFunc) gosh.HandlerFunc {
	return func(c *gosh.Context) error {
		// 检查是否跳过中间件或请求体为空
		if SkipHandler(c, skippers...) || c.Request.Body == nil {
			c.Next()
			return nil
		}

		var requestBody []byte
		isGzip := false
		safe := &io.LimitedReader{R: c.Request.Body, N: GetMaxMemory()}

		// 检查请求是否使用gzip压缩
		if c.Header("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(safe)
			if err == nil {
				isGzip = true
				requestBody, _ = io.ReadAll(reader)
			}
		}

		// 如果不是gzip压缩或解压缩时出错，则直接读取请求体
		if !isGzip {
			requestBody, _ = io.ReadAll(safe)
		}

		// 关闭原始请求体，创建新的MaxBytesReader，将复制的请求体作为缓冲区
		c.Request.Body.Close()
		bf := bytes.NewBuffer(requestBody)
		c.Request.Body = http.MaxBytesReader(c.ResponseWriter, io.NopCloser(bf), GetMaxMemory())
		c.SetContextValue(reqBodyKey, requestBody)
		c.Next()
		return nil
	}
}
