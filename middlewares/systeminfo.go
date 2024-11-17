/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2023-07-28 00:50:58
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 15:06:05
 * @FilePath: \gosh\middlewares\systeminfo.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/moment"
	"github.com/kamalyes/go-toolbox/pkg/osx"
)

// SystemInfo 存储系统信息
type SystemInfo struct {
	ServerName   string // 服务器名称
	Runtime      string // 运行时间
	GoroutineNum string // goroutine数量
	CPUNum       string // cpu核数
	UsedMem      string // 当前内存使用量
	TotalMem     string // 总分配的内存
	SysMem       string // 系统内存占用量
	Lookups      string // 指针查找次数
	Mallocs      string // 内存分配次数
	Frees        string // 内存释放次数
	LastGCTime   string // 距离上次GC时间
	NextGC       string // 下次GC内存回收量
	PauseTotalNs string // GC暂停时间总量
	PauseNs      string // 上次GC暂停时间
	HeapInuse    string // 正在使用的堆内存
}

// SizeFormat 格式化文件大小
func SizeFormat(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	humanaBytes := func(s uint64, base float64, sizes []string) string {
		if s < 10 {
			return fmt.Sprintf("%d B", s)
		}
		e := math.Floor(logSize(float64(s), base))
		suffix := sizes[int(e)]
		val := float64(s) / math.Pow(base, math.Floor(e))
		f := "%.0f"
		if val < 10 {
			f = "%.1f"
		}
		return fmt.Sprintf(f+" %s", val, suffix)
	}
	return humanaBytes(s, 1024, sizes)
}

func logSize(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

// NewSystemInfo 创建 SystemInfo 结构体的实例
func NewSystemInfo() *SystemInfo {
	var afterLastGC string
	mstat := &runtime.MemStats{}
	runtime.ReadMemStats(mstat)
	// 记录当前时间
	_, costTime, _ := moment.GetCurrentTimeInfo()
	if mstat.LastGC != 0 {
		afterLastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(mstat.LastGC))/1000/1000/1000)
	} else {
		afterLastGC = "0"
	}

	return &SystemInfo{
		ServerName:   osx.SafeGetHostName(),                                                                                      // 获取服务器名称
		Runtime:      fmt.Sprintf("%d天%d小时%d分%d秒", costTime/(3600*24), costTime%(3600*24)/3600, costTime%3600/60, costTime%(60)), // 计算运行时间
		GoroutineNum: strconv.Itoa(runtime.NumGoroutine()),                                                                       // 获取goroutine数量
		CPUNum:       strconv.Itoa(runtime.NumCPU()),                                                                             // 获取cpu核数
		HeapInuse:    SizeFormat(uint64(mstat.HeapInuse)),                                                                        // 正在使用的堆内存
		UsedMem:      SizeFormat(uint64(mstat.Alloc)),                                                                            // 当前内存使用量
		TotalMem:     SizeFormat(uint64(mstat.TotalAlloc)),                                                                       // 总分配的内存
		SysMem:       SizeFormat(uint64(mstat.Sys)),                                                                              // 系统内存占用量
		Lookups:      strconv.FormatUint(mstat.Lookups, 10),                                                                      // 指针查找次数
		Mallocs:      strconv.FormatUint(mstat.Mallocs, 10),                                                                      // 内存分配次数
		Frees:        strconv.FormatUint(mstat.Frees, 10),                                                                        // 内存释放次数
		LastGCTime:   afterLastGC,                                                                                                // 距离上次GC时间
		NextGC:       SizeFormat(uint64(mstat.NextGC)),                                                                           // 下次GC内存回收量
		PauseTotalNs: fmt.Sprintf("%.3fs", float64(mstat.PauseTotalNs)/1000/1000/1000),                                           // GC暂停时间总量
		PauseNs:      fmt.Sprintf("%.3fs", float64(mstat.PauseNs[(mstat.NumGC+255)%256])/1000/1000/1000),                         // 上次GC暂停时间
	}
}
