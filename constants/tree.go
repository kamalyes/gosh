/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 21:19:12
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 00:55:29
 * @FilePath: \gosh\constants\tree.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package constants

// 定义字符串类型的路径参数前缀、通配符和路径分隔符
const (
	PathSeparator   = '/' // 定义路径分隔符
	PathParamPrefix = ':' // 定义路径参数前缀
	WildcardSymbol  = '*' // 定义通配符
)

var (
	PathParamPrefixStr = ":"        // 路径参数的前缀
	PathSeparatorStr   = "/"        // 路径分隔符
	WildcardSymbolStr  = "*"        // 通配符
	IllegalPath        = ":*?\"<>|" // 非法字符
)

// 定义常量来表示特定的路径值
const (
	CurrentDirectory = "." // 当前目录
	RootDirectory    = "/" // 根目录
)
