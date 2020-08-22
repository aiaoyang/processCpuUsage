package yssender

import (
	"time"

	influxc "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
)

const meansure = "system"

// MySender 自定义Sender
type MySender struct {
	API api.WriteAPI
}

// Send 实现 customMetric Sender 的Send方法
func (s *MySender) Send(k map[string]string, v map[string]interface{}) {
	point := influxc.NewPoint(
		meansure,
		k,
		v,
		time.Now(),
	)
	s.API.WritePoint(point)
}
