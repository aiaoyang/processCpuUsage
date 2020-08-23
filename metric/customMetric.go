package metric

import (
	"context"
	"fmt"
	"log"
	"time"
)

// CustomMetric CustomMetric
type CustomMetric struct {
	// 初始化后不可修改
	sender Sender
	Tag
	Metric
}

// NewCustomMetric 初始化自定义指标
func NewCustomMetric(sender Sender) *CustomMetric {
	return &CustomMetric{
		sender: sender,
		Tag:    make(Tag),
		Metric: make(Metric),
	}
}

// Send 发送结构体内的数据指标
func (c *CustomMetric) Send() {
	if c.sender == nil {
		log.Fatal(fmt.Errorf("metric sender is nil"))
	}
	tag := c.Tag.CopyToMap()
	metric := c.Metric.CopyToMap()
	c.FlushMetric()
	c.sender.Send(tag, metric)
}

// MetricConsumer MetricConsumer
func (c *CustomMetric) MetricConsumer(ctx context.Context, reciver <-chan Metric) {
	for {
		select {
		case <-ctx.Done():
			return

		// 接收到指标参数，则将其推送至influxdb
		case res := <-reciver:

			// 合并指标
			c.Metric.Add(res)

			// 发送数据
			c.Send()

		default:

			time.Sleep(time.Second * 3)

		}
	}
}

// Flush 清空保存的内容
func (c *CustomMetric) Flush() {
	c.FlushMetric()
	c.FlushTag()
}

// FlushTag 清空tag
func (c *CustomMetric) FlushTag() {
	for k := range c.Tag {
		delete(c.Tag, k)
	}
}

// FlushMetric 清空metric
func (c *CustomMetric) FlushMetric() {
	for k := range c.Metric {
		delete(c.Metric, k)
	}
}
