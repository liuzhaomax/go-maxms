package core

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func InitPrometheusRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		// 用于收集与 Go 运行时相关的指标数据，例如 Goroutine 数量、内存使用情况、GC（垃圾回收）活动等
		collectors.NewGoCollector(),
		// 用于收集与进程相关的指标数据，例如进程的 CPU 使用情况、内存使用情况等
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	return registry
}
