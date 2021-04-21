package configs

// AliYunConfig 访问 ·阿里云配置中心· 认证配置
type AliYunConfig struct {
	Endpoint    string `yaml:"endpoint"`
	NamespaceID string `yaml:"namespaceID"`
	AccessKey   string `yaml:"accessKey"`
	SecretKey   string `yaml:"secretKey"`

	DataID string `yaml:"dataID"`
	Group  string `yaml:"group"`
}
