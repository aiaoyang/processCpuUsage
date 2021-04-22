package monitors

import "github.com/aiaoyang/processCpuUsage/pkg/metric"

type Monitor interface {
	Mon(sender metric.Sender)
}
