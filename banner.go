/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-21 17:30:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 17:55:15
 * @FilePath: \gosh\banner.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"fmt"
	"os"
	"sync"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/kamalyes/gosh/constants"
)

// GetAppName 获取应用程序名称，如果设置了，则使用环境变量
func GetAppName() string {
	return os.Getenv(constants.DefaultAppPathVariableName)
}

// BannerConfig 用于存储 Banner 配置
type BannerConfig struct {
	banner   string
	title    string
	subtitle string
	mu       sync.RWMutex
}

// NewBannerConfig 创建一个新的 BannerConfig 实例
func NewBannerConfig() *BannerConfig {
	return &BannerConfig{
		banner:   constants.DefaultBanner,
		title:    constants.DefaultTitle,
		subtitle: constants.DefaultSubtitle,
	}
}

// SetBanner 设置 Banner
func (bc *BannerConfig) SetBanner(banner string) *BannerConfig {
	return syncx.WithLockReturnValue(&bc.mu, func() *BannerConfig {
		bc.banner = stringx.Trim(banner)
		return bc
	})
}

// SetTitle 设置 Title
func (bc *BannerConfig) SetTitle(title string) *BannerConfig {
	return syncx.WithLockReturnValue(&bc.mu, func() *BannerConfig {
		bc.title = stringx.Trim(title)
		return bc
	})
}

// SetSubtitle 设置 Subtitle
func (bc *BannerConfig) SetSubtitle(subtitle string) *BannerConfig {
	return syncx.WithLockReturnValue(&bc.mu, func() *BannerConfig {
		bc.subtitle = stringx.Trim(subtitle)
		return bc
	})
}

// GetBanner 获取 Banner
func (bc *BannerConfig) GetBanner() string {
	return syncx.WithRLockReturnValue(&bc.mu, func() string {
		return bc.banner
	})
}

// GetTitle 获取 Title
func (bc *BannerConfig) GetTitle() string {
	return syncx.WithRLockReturnValue(&bc.mu, func() string {
		return bc.title
	})
}

// GetSubtitle 获取 Subtitle
func (bc *BannerConfig) GetSubtitle() string {
	return syncx.WithRLockReturnValue(&bc.mu, func() string {
		return bc.subtitle
	})
}

// Print 打印可配置的 Print
func (bc *BannerConfig) Print() {
	syncx.WithRLockReturn(&bc.mu, func() (string, error) {
		banner := fmt.Sprintf("%s\n%s\n%s\n", bc.banner, bc.title, bc.subtitle)
		fmt.Println(banner)
		return banner, nil
	})
}
