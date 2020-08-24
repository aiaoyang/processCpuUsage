package metric

import (
	"reflect"

	"github.com/aiaoyang/processCpuUsage/sysusage"
)

type MetricType string

const (
	// LOAD1 1分钟负载
	LOAD1 MetricType = "load1"

	// LOAD5 5分钟负载
	LOAD5 MetricType = "load5"

	// LOAD15 15分钟负载
	LOAD15 MetricType = "load15"

	// SYSCPU 系统cpu使用率
	SYSCPU MetricType = "sys_cpu"

	// SYSMEM 系统内存使用率
	SYSMEM MetricType = "sys_mem"

	// PROCESS_CPU 进程cpu使用率
	PROCESS_CPU MetricType = "process_cpu"

	// PROCESS_MEM 进程内存使用率
	PROCESS_MEM MetricType = "process_mem"

	// FD 进程文件打开数
	FD MetricType = "fd"

	// PID 进程pid
	PID = "pid"
)

// Metric 重新封装map指标
type Metric map[MetricType]interface{}

// NewMetric 初始化一个metric
func NewMetric() Metric {
	return make(Metric)
}

// Copy 值拷贝
func (m Metric) Copy() Metric {
	tmp := make(Metric)
	for k, v := range m {
		tmp[k] = v
	}
	return tmp
}

// CopyToMap 值拷贝
func (m Metric) CopyToMap() map[string]interface{} {
	tmp := make(map[string]interface{})
	for k, v := range m {
		tmp[string(k)] = v
	}
	return tmp
}

// Insert 添加键值对
func (m Metric) Insert(k MetricType, v interface{}) {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Float64:
		m[k] = float64(v.(sysusage.Usage))
	case reflect.Int:
		m[k] = v.(int)
	default:
		m[k] = -1
	}
}

// Add 合并指标
func (m Metric) Add(subs ...Metric) {
	for i := 0; i < len(subs); i++ {
		for k, v := range subs[i] {
			m[k] = v
		}
	}
}

// AddMap 合并指标
func (m Metric) AddMap(subs ...map[string]interface{}) {
	for i := 0; i < len(subs); i++ {
		for k, v := range subs[i] {
			m[MetricType(k)] = v
		}
	}
}
