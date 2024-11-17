/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-15 23:26:10
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 16:39:16
 * @FilePath: \gosh\status_code.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"sync"
)

type StatusCodeInterface interface {
	SetStatusCode(code StatusCode, msg string)
	GetStatusCodeMsg(code StatusCode) string
	InitStatusCodes(customCodes map[StatusCode]string)
}

type StatusCode int

const (
	StatusContinue           StatusCode = 100 // 继续
	StatusSwitchingProtocols StatusCode = 101 // 交换协议
	StatusProcessing         StatusCode = 102 // 处理中
	StatusEarlyHints         StatusCode = 103 // 提示信息

	StatusOK                   StatusCode = 200 // 成功
	StatusCreated              StatusCode = 201 // 已创建
	StatusAccepted             StatusCode = 202 // 已接受
	StatusNonAuthoritativeInfo StatusCode = 203 // 非权威信息
	StatusNoContent            StatusCode = 204 // 无内容
	StatusResetContent         StatusCode = 205 // 重置内容
	StatusPartialContent       StatusCode = 206 // 部分内容
	StatusMultiStatus          StatusCode = 207 // 多状态
	StatusAlreadyReported      StatusCode = 208 // 已报告
	StatusIMUsed               StatusCode = 226 // 使用了IM

	StatusMultipleChoices   StatusCode = 300 // 多种选择
	StatusMovedPermanently  StatusCode = 301 // 永久移动
	StatusFound             StatusCode = 302 // 已找到
	StatusSeeOther          StatusCode = 303 // 另请参见
	StatusNotModified       StatusCode = 304 // 未修改
	StatusUseProxy          StatusCode = 305 // 使用代理
	StatusTemporaryRedirect StatusCode = 307 // 临时重定向
	StatusPermanentRedirect StatusCode = 308 // 永久重定向

	StatusBadRequest                   StatusCode = 400 // 错误请求
	StatusUnauthorized                 StatusCode = 401 // 未授权
	StatusPaymentRequired              StatusCode = 402 // 需要付款
	StatusForbidden                    StatusCode = 403 // 禁止访问
	StatusNotFound                     StatusCode = 404 // 未找到
	StatusMethodNotAllowed             StatusCode = 405 // 方法不允许
	StatusNotAcceptable                StatusCode = 406 // 不可接受
	StatusProxyAuthRequired            StatusCode = 407 // 需要代理授权
	StatusRequestTimeout               StatusCode = 408 // 请求超时
	StatusConflict                     StatusCode = 409 // 冲突
	StatusGone                         StatusCode = 410 // 已删除
	StatusLengthRequired               StatusCode = 411 // 长度必需
	StatusPreconditionFailed           StatusCode = 412 // 先决条件失败
	StatusRequestEntityTooLarge        StatusCode = 413 // 请求实体过大
	StatusRequestURITooLong            StatusCode = 414 // 请求URI过长
	StatusUnsupportedMediaType         StatusCode = 415 // 不支持的媒体类型
	StatusRequestedRangeNotSatisfiable StatusCode = 416 // 请求范围不符合要求
	StatusExpectationFailed            StatusCode = 417 // 预期失败
	StatusTeapot                       StatusCode = 418 // 我是茶壶
	StatusMisdirectedRequest           StatusCode = 421 // 误导的请求
	StatusUnprocessableEntity          StatusCode = 422 // 无法处理的实体
	StatusLocked                       StatusCode = 423 // 已锁定
	StatusFailedDependency             StatusCode = 424 // 依赖失败
	StatusTooEarly                     StatusCode = 425 // 太早
	StatusUpgradeRequired              StatusCode = 426 // 需要升级
	StatusPreconditionRequired         StatusCode = 428 // 先决条件要求
	StatusTooManyRequests              StatusCode = 429 // 请求过多
	StatusRequestHeaderFieldsTooLarge  StatusCode = 431 // 请求头字段太大
	StatusUnavailableForLegalReasons   StatusCode = 451 // 由于法律原因不可用

	StatusInternalServerError           StatusCode = 500 // 服务器内部错误
	StatusNotImplemented                StatusCode = 501 // 未实现
	StatusBadGateway                    StatusCode = 502 // 网关错误
	StatusServiceUnavailable            StatusCode = 503 // 服务不可用
	StatusGatewayTimeout                StatusCode = 504 // 网关超时
	StatusHTTPVersionNotSupported       StatusCode = 505 // 不支持的HTTP版本
	StatusVariantAlsoNegotiates         StatusCode = 506 // 可协商的变体
	StatusInsufficientStorage           StatusCode = 507 // 存储空间不足
	StatusLoopDetected                  StatusCode = 508 // 检测到循环
	StatusNotExtended                   StatusCode = 510 // 未扩展
	StatusNetworkAuthenticationRequired StatusCode = 511 // 需要网络认证
)

// StatusCodeMsgMap 用于存储状态码Key映射关系
var StatusCodeMsgMap = struct {
	sync.RWMutex
	mapping map[StatusCode]string
}{
	mapping: map[StatusCode]string{
		StatusContinue:                      "Continue",
		StatusSwitchingProtocols:            "Switching Protocols",
		StatusProcessing:                    "Processing",
		StatusEarlyHints:                    "Early Hints",
		StatusOK:                            "OK",
		StatusCreated:                       "Created",
		StatusAccepted:                      "Accepted",
		StatusNonAuthoritativeInfo:          "Non-Authoritative Information",
		StatusNoContent:                     "No Content",
		StatusResetContent:                  "Reset Content",
		StatusPartialContent:                "Partial Content",
		StatusMultiStatus:                   "Multi-Status",
		StatusAlreadyReported:               "Already Reported",
		StatusIMUsed:                        "IM Used",
		StatusMultipleChoices:               "Multiple Choices",
		StatusMovedPermanently:              "Moved Permanently",
		StatusFound:                         "Found",
		StatusSeeOther:                      "See Other",
		StatusNotModified:                   "Not Modified",
		StatusUseProxy:                      "Use Proxy",
		StatusTemporaryRedirect:             "Temporary Redirect",
		StatusPermanentRedirect:             "Permanent Redirect",
		StatusBadRequest:                    "Bad Request",
		StatusUnauthorized:                  "Unauthorized",
		StatusPaymentRequired:               "Payment Required",
		StatusForbidden:                     "Forbidden",
		StatusNotFound:                      "Not Found",
		StatusMethodNotAllowed:              "Method Not Allowed",
		StatusNotAcceptable:                 "Not Acceptable",
		StatusProxyAuthRequired:             "Proxy Authentication Required",
		StatusRequestTimeout:                "Request Timeout",
		StatusConflict:                      "Conflict",
		StatusGone:                          "Gone",
		StatusLengthRequired:                "Length Required",
		StatusPreconditionFailed:            "Precondition Failed",
		StatusRequestEntityTooLarge:         "Request Entity Too Large",
		StatusRequestURITooLong:             "Request URI Too Long",
		StatusUnsupportedMediaType:          "Unsupported Media Type",
		StatusRequestedRangeNotSatisfiable:  "Requested Range Not Satisfiable",
		StatusExpectationFailed:             "Expectation Failed",
		StatusTeapot:                        "I'm a teapot",
		StatusMisdirectedRequest:            "Misdirected Request",
		StatusUnprocessableEntity:           "Unprocessable Entity",
		StatusLocked:                        "Locked",
		StatusFailedDependency:              "Failed Dependency",
		StatusTooEarly:                      "Too Early",
		StatusUpgradeRequired:               "Upgrade Required",
		StatusPreconditionRequired:          "Precondition Required",
		StatusTooManyRequests:               "Too Many Requests",
		StatusRequestHeaderFieldsTooLarge:   "Request Header Fields Too Large",
		StatusUnavailableForLegalReasons:    "Unavailable For Legal Reasons",
		StatusInternalServerError:           "Internal Server Error",
		StatusNotImplemented:                "Not Implemented",
		StatusBadGateway:                    "Bad Gateway",
		StatusServiceUnavailable:            "Service Unavailable",
		StatusGatewayTimeout:                "Gateway Timeout",
		StatusHTTPVersionNotSupported:       "HTTP Version Not Supported",
		StatusVariantAlsoNegotiates:         "Variant Also Negotiates",
		StatusInsufficientStorage:           "Insufficient Storage",
		StatusLoopDetected:                  "Loop Detected",
		StatusNotExtended:                   "Not Extended",
		StatusNetworkAuthenticationRequired: "Network Authentication Required",
	},
}

// SetStatusCode 设置状态码
func SetStatusCode(code StatusCode, msg string) {
	StatusCodeMsgMap.Lock()
	defer StatusCodeMsgMap.Unlock()
	StatusCodeMsgMap.mapping[code] = msg
}

// GetStatusCodeText 根据状态码获取对应的错误消息
func GetStatusCodeText(code StatusCode) string {
	StatusCodeMsgMap.RLock()
	defer StatusCodeMsgMap.RUnlock()
	return StatusCodeMsgMap.mapping[code]
}

// InitStatusCodes 初始化状态码和消息
func InitStatusCodes(customCodes map[StatusCode]string) {
	StatusCodeMsgMap.Lock()
	defer StatusCodeMsgMap.Unlock()

	for code, msg := range customCodes {
		StatusCodeMsgMap.mapping[code] = msg
	}
}
