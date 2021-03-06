package db

import (
	"time"

	"github.com/aiaoyang/processCpuUsage/pkg/metric"
	influxc "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

const meansure = "system"

// InfluxDBSender 自定义Sender
type InfluxDBSender struct {
	API api.WriteAPI
}

// NewInfluxDBSender 初始化influxdb类型的sender
func NewInfluxDBSender(api api.WriteAPI) metric.Sender {
	return &InfluxDBSender{
		API: api,
	}
}

// Send 实现 customMetric Sender 的Send方法
func (s *InfluxDBSender) Send(tag map[string]string, value map[string]interface{}) {
	point := influxc.NewPoint(
		meansure,
		tag,
		value,
		time.Now(),
	)
	s.API.WritePoint(point)
}
