package metric

// Sender 数据传输发送接口
type Sender interface {
	// ReNew() bool
	// Health() bool
	Send(tag map[string]string, metric map[string]interface{})
}
