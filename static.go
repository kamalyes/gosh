/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-15 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 09:24:36
 * @FilePath: \gosh\static.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kamalyes/gosh/constants"
	"github.com/kamalyes/gosh/errorsx"
)

// Static 注册一个指向服务端本地目录的静态路由
// relativePath: 路由的相对路径
// localPath: 本地文件系统中的路径
// listDir: 是否允许列出目录内容
func (group *RouterGroup) Static(relativePath, localPath string, listDir bool) {
	// 验证相对路径和本地路径的有效性
	absLocalPath, err := validatePaths(relativePath, localPath)
	if err != nil {
		panic(err) // 如果路径无效，抛出异常
	}

	// 定义处理静态文件请求的处理函数
	handler := func(ctx *Context) error {
		return serveStaticFile(ctx, absLocalPath, listDir, group.Engine)
	}

	// 在路由组中注册 GET 和 HEAD 方法的处理函数
	// 注意：这里使用了命名参数而不是通配符，但通配符仍然可以用于更复杂的匹配
	// 如果需要完全匹配文件路径，可以考虑不使用命名参数，而是让handler内部处理
	// 但为了示例，这里保留命名参数并添加额外的逻辑来处理目录
	finalURLPath := path.Join(relativePath, "/*filepath/") // 这里的通配符仍然保留，但需要在handler中处理
	group.GET(finalURLPath, handler)
	group.HEAD(finalURLPath, handler)
}

// validatePaths 检查相对路径和本地路径的有效性，并返回本地路径的绝对路径。
// - relativePath: 路由的相对路径（用于验证是否包含非法字符）。
// - localPath: 本地文件系统中的路径（将转换为绝对路径）。
// 返回绝对路径和可能的错误。
func validatePaths(relativePath, localPath string) (absLocalPath string, err error) {
	// 检查相对路径是否包含非法字符
	if strings.ContainsAny(relativePath, constants.IllegalPath) {
		return "", fmt.Errorf("相对路径不能使用非法字符: %s", relativePath)
	}

	// 清理本地路径
	cleanedLocalPath := filepath.Clean(localPath)
	if cleanedLocalPath == constants.CurrentDirectory || cleanedLocalPath == constants.RootDirectory {
		// 防止清理后的路径是根目录或当前目录，这可能不是预期的行为
		return "", fmt.Errorf("本地路径不能是根目录或当前目录: %s", localPath)
	}

	// 获取本地路径的绝对路径
	absLocalPath, err = filepath.Abs(cleanedLocalPath)
	if err != nil {
		return "", fmt.Errorf("无法获取本地路径的绝对路径: %v", err)
	}

	// 检查绝对路径是否为空（虽然filepath.Abs通常不会返回空字符串，但这里作为额外的检查）
	if absLocalPath == "" {
		return "", fmt.Errorf("获取的绝对路径为空（这通常不应该发生）: %s", localPath)
	}

	// 还可以添加其他检查，比如验证路径是否指向一个有效的目录或文件

	return absLocalPath, nil
}

// serveStaticFile 处理静态文件请求
// ctx: 请求上下文
// filePath: 本地文件的完整路径（包括请求的文件名或目录）
// listDir: 是否允许列出目录内容
func serveStaticFile(ctx *Context, filePath string, listDir bool, engine *Engine) error {
	// 检查文件是否存在和其他属性
	fileInfo, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		// 尝试修复路径，如果请求的是目录且末尾没有斜杠，则重定向到带斜杠的URL（可选）
		if fileInfoDir, dirErr := os.Stat(filePath + constants.PathSeparatorStr); !os.IsNotExist(dirErr) && dirErr == nil && fileInfoDir.IsDir() {
			return handleError(ctx, engine, errorsx.ErrFileNotFound, http.StatusNotFound)
		}
		return handleError(ctx, engine, errorsx.ErrFileNotFound, http.StatusNotFound)
	}
	if err != nil {
		if os.IsPermission(err) {
			return handleError(ctx, engine, errorsx.ErrAccessDenied, http.StatusForbidden)
		}
		return handleError(ctx, engine, errorsx.ErrInternalServerError, http.StatusInternalServerError)
	}

	// 检查是否为目录且是否允许列出目录内容
	if fileInfo.IsDir() {
		if !listDir {
			// 临时返回403，应替换为实际列出目录的逻辑
			return handleError(ctx, engine, errorsx.ErrDirectoryAccessForbidden, http.StatusForbidden)
		}
	}

	// 处理文件请求
	http.ServeFile(ctx.ResponseWriter, ctx.Request, filePath)
	return nil
}
