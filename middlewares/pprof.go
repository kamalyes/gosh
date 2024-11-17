/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2024-08-01 13:50:19
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2024-11-17 22:55:15
 * @FilePath: \gosh\middlewares\pprof.go
 * @Description:
 *
 * Copyright (c) 2024 by kamalyes, All Rights Reserved.
 */
package middlewares

import (
	"fmt"
	"net/http/pprof"

	"github.com/kamalyes/gosh"
)

const (
	// pprof的默认url前缀
	DefaultPrefix = "/debug/pprof"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

func PprofHandler(c *gosh.Context) error {
	m := NewSystemInfo()
	// 定义一个常量，用于格式化输出
	const formatString = "%s:%s\n"
	info := fmt.Sprintf(formatString, "服务器", m.ServerName)
	info += fmt.Sprintf(formatString, "运行时间", m.Runtime)
	info += fmt.Sprintf(formatString, "goroutine数量", m.GoroutineNum)
	info += fmt.Sprintf(formatString, "CPU核数", m.CPUNum)
	info += fmt.Sprintf(formatString, "当前内存使用量", m.UsedMem)
	info += fmt.Sprintf(formatString, "当前堆内存使用量", m.HeapInuse)
	info += fmt.Sprintf(formatString, "总分配的内存", m.TotalMem)
	info += fmt.Sprintf(formatString, "系统内存占用量", m.SysMem)
	info += fmt.Sprintf(formatString, "指针查找次数", m.Lookups)
	info += fmt.Sprintf(formatString, "内存分配次数", m.Mallocs)
	info += fmt.Sprintf(formatString, "内存释放次数", m.Frees)
	info += fmt.Sprintf(formatString, "距离上次GC时间", m.LastGCTime)
	info += fmt.Sprintf(formatString, "下次GC内存回收量", m.NextGC)
	info += fmt.Sprintf(formatString, "GC暂停时间总量", m.PauseTotalNs)
	info += fmt.Sprintf(formatString, "上次GC暂停时间", m.PauseNs)
	_, _ = fmt.Fprint(c.ResponseWriter, info)
	return nil
}

// 从. net/http/pprof包注册标准HandlerFuncs
// 提供的gosh.Router。prefixOptions是可选的。如果不是prefixOptions，
// 使用默认的路径前缀，否则第一个prefixOptions将是路径前缀。
func Register(r gosh.Router, prefixOptions ...string) {
	PprofRouteRegister(r, prefixOptions...)
}

// 将标准HandlerFuncs从net/http/pprof包中注册到
// 使用默认的路径前缀，否则第一个prefixOptions将是路径前缀。
func PprofRouteRegister(rg gosh.Router, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)

	prefixRouter := rg.Group(prefix)
	{
		prefixRouter.GET("/", gosh.WrapF(pprof.Index))
		prefixRouter.GET("/sysinfo", PprofHandler)
		prefixRouter.GET("/cmdline", gosh.WrapF(pprof.Cmdline))
		prefixRouter.GET("/profile", gosh.WrapF(pprof.Profile))
		prefixRouter.POST("/symbol", gosh.WrapF(pprof.Symbol))
		prefixRouter.GET("/symbol", gosh.WrapF(pprof.Symbol))
		prefixRouter.GET("/trace", gosh.WrapF(pprof.Trace))
		prefixRouter.GET("/allocs", gosh.WrapH(pprof.Handler("allocs")))
		prefixRouter.GET("/block", gosh.WrapH(pprof.Handler("block")))
		prefixRouter.GET("/goroutine", gosh.WrapH(pprof.Handler("goroutine")))
		prefixRouter.GET("/heap", gosh.WrapH(pprof.Handler("heap")))
		prefixRouter.GET("/mutex", gosh.WrapH(pprof.Handler("mutex")))
		prefixRouter.GET("/threadcreate", gosh.WrapH(pprof.Handler("threadcreate")))
	}
}
