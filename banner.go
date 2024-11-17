/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-16 03:00:16
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-16 03:20:09
 * @FilePath: \gosh\banner.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	"fmt"

	"github.com/kamalyes/gosh/constants"
)

// BannerConfig 用于存储 Banner 配置
type BannerConfig struct {
	Banner   string
	Title    string
	Subtitle string
}

// PrintBanner 打印可配置的 Banner
func PrintBanner(config BannerConfig) {
	banner := fmt.Sprintf("%s\n   %s\n   %s\n", config.Banner, config.Title, config.Subtitle)
	fmt.Println(banner)
}

func DefaultBannerConfig() *BannerConfig {
	return &BannerConfig{
		Banner:   constants.DefaultBanner,
		Title:    constants.DefaultTitle,
		Subtitle: constants.DefaultSubtitle,
	}
}
