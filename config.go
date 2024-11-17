/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 13:55:55
 * @FilePath: \gosh\config.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package gosh

import (
	goconfig "github.com/kamalyes/go-config"
)

// setDefaultConfig 设置默认配置
func setDefaultConfig() Config {
	config := Config{
		MaxMultipartMemory:     defaultMaxMultipartMemory,
		HandleMethodNotAllowed: false,
		AppBanner:              DefaultBannerConfig(),
		KmSingleConfig: &goconfig.SingleConfig{
			Zap: DefaultKmZipConfig(),
		},
	}
	return config
}

func mergeDefaultConfig(defaultConfig, customConfig Config) Config {
	if customConfig.AppBanner == nil {
		defaultConfig.AppBanner = DefaultBannerConfig()
	}

	if customConfig.KmSingleConfig == nil {
		defaultConfig.KmSingleConfig.Zap = DefaultKmZipConfig()
	}
	if customConfig.AppBanner.Banner == "" {
		defaultConfig.AppBanner.Banner = DefaultBannerConfig().Banner
	}

	if customConfig.AppBanner.Subtitle == "" {
		defaultConfig.AppBanner.Subtitle = DefaultBannerConfig().Subtitle
	}

	if customConfig.AppBanner.Title == "" {
		defaultConfig.AppBanner.Title = DefaultBannerConfig().Title
	}

	if customConfig.AppName == "" {
		defaultConfig.AppName = customConfig.AppName
	}

	if customConfig.AppName != "" {
		defaultConfig.AppName = customConfig.AppName
	}

	return defaultConfig
}
