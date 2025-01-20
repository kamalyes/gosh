/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 21:19:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 10:15:55
 * @FilePath: \gosh\constants\headers.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package constants

// HTTP 头部相关关键字常量
const (
	HeaderContentTypeKey     = "Content-Type"
	HeaderLocationKey        = "Location"
	HeaderContentLengthKey   = "Content-Length"
	HeaderContentEncodingKey = "Content-Encoding"
	HeaderOriginKey          = "Origin"
)

// ContentEncoding 相关的常量
const (
	ContentEncodingGzip = "gzip"
)

// CORS 相关的常量
const (
	AccessControlAllowOriginKey      = "Access-Control-Allow-Origin"
	AccessControlAllowMethodsKey     = "Access-Control-Allow-Methods"
	AccessControlAllowHeadersKey     = "Access-Control-Allow-Headers"
	AccessControlExposeHeadersKey    = "Access-Control-Expose-Headers"
	AccessControlMaxAgeKey           = "Access-Control-Max-Age"
	AccessControlAllowCredentialsKey = "Access-Control-Allow-Credentials"
)
