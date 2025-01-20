/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-11-11 08:59:07
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-01-20 17:56:01
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
		AppBanner:              NewBannerConfig(),
		KmSingleConfig: &goconfig.SingleConfig{
			Zap: DefaultKmZipConfig(),
		},
	}
	return config
}

func mergeDefaultConfig(defaultConfig, customConfig Config) Config {
	if customConfig.MaxMultipartMemory != 0 {
		defaultConfig.MaxMultipartMemory = customConfig.MaxMultipartMemory
	}

	if customConfig.Recovery {
		defaultConfig.Recovery = customConfig.Recovery
	}

	if customConfig.HandleMethodNotAllowed {
		defaultConfig.HandleMethodNotAllowed = customConfig.HandleMethodNotAllowed
	}

	if customConfig.BeforeHandler != nil {
		defaultConfig.BeforeHandler = customConfig.BeforeHandler
	}

	if customConfig.ErrorHandler != nil {
		defaultConfig.ErrorHandler = customConfig.ErrorHandler
	}

	if customConfig.AfterHandler != nil {
		defaultConfig.AfterHandler = customConfig.AfterHandler
	}

	defaultBannerConfig := NewBannerConfig()

	if customConfig.AppBanner == nil {
		defaultConfig.AppBanner = defaultBannerConfig
	}

	if customConfig.AppBanner != nil {
		if customConfig.AppBanner.banner == "" {
			defaultConfig.AppBanner.SetBanner(defaultBannerConfig.GetBanner())
		}

		if customConfig.AppBanner.subtitle == "" {
			defaultConfig.AppBanner.SetSubtitle(defaultBannerConfig.GetSubtitle())
		}

		if customConfig.AppBanner.title == "" {
			defaultConfig.AppBanner.SetTitle(defaultBannerConfig.GetTitle())
		}
	}

	if customConfig.KmSingleConfig == nil {
		defaultConfig.KmSingleConfig.Zap = DefaultKmZipConfig()
	}

	if customConfig.AppName != "" {
		defaultConfig.AppName = customConfig.AppName
	}

	return defaultConfig
}
