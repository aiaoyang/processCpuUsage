package metric

// Tag influxdb 标签
type Tag map[TagType]string

// TagType 标签类型
type TagType string

const (
	// HOSTNAME 主机名
	HOSTNAME TagType = "hostname"

	// ENV 主机环境属性，测试或正式或提神
	ENV TagType = "env"
)

// NewTag 初始化一个 Tag
func NewTag() Tag {
	return make(Tag)
}

// Copy 值拷贝
func (m Tag) Copy() Tag {
	tmp := make(Tag)
	for k, v := range m {
		tmp[k] = v
	}
	return tmp
}

// CopyToMap 值拷贝
func (m Tag) CopyToMap() map[string]string {
	tmp := make(map[string]string)
	for k, v := range m {
		tmp[string(k)] = v
	}
	return tmp
}

// Insert 添加键值对
func (m Tag) Insert(k TagType, v string) {
	m[k] = v
}

// Add 合并指标
func (m Tag) Add(subs ...Tag) {
	for i := 0; i < len(subs); i++ {
		for k, v := range subs[i] {
			m[k] = v
		}
	}
}

// AddMap 合并指标
func (m Tag) AddMap(subs ...map[string]string) {
	for i := 0; i < len(subs); i++ {
		for k, v := range subs[i] {
			m[TagType(k)] = v
		}
	}
}
