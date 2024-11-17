/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-15 21:51:51
 * @FilePath: \gosh\static_file.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"net/http"
	"strings"

	"github.com/kamalyes/gosh/constants"
)

// StaticFile 注册一个指向服务端本地文件的静态路由
func (group *RouterGroup) StaticFile(relativePath, filePath string) {
	group.staticFileHandler(relativePath, func(ctx *Context) error {
		ctx.ServeFile(filePath)
		return nil
	})
}

// StaticFileFS 与StaticFile函数类型，但可以自定义文件系统
func (group *RouterGroup) StaticFileFS(relativePath, filePath string, fs http.FileSystem) {
	group.staticFileHandler(relativePath, func(ctx *Context) error {
		ctx.FileFromFS(filePath, fs)
		return nil
	})
}

func (group *RouterGroup) staticFileHandler(relativePath string, handler HandlerFunc) {
	if strings.ContainsAny(relativePath, constants.IllegalPath) {
		panic("URL parameters can not be used when serving a staticNode file")
	}
	group.GET(relativePath, handler)
	group.HEAD(relativePath, handler)
}
