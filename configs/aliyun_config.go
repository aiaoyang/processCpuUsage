package configs

// AliyunAcm 访问 ·阿里云配置中心· 认证配置
type AliyunAcm struct {
	Endpoint    string `yaml:"endpoint"`
	NamespaceID string `yaml:"namespaceID"`
	AccessKey   string `yaml:"accessKey"`
	SecretKey   string `yaml:"secretKey"`

	DataID string `yaml:"dataID"`
	Group  string `yaml:"group"`
}

type AcmInnerConfig struct {
	Hostname     string   `yaml:"hostName"`
	StoredPids   []string `yaml:"storedPids"`
	ProcessRegex string   `yaml:"processRegex"`
	IsMonitorOn  bool     `yaml:"isMonitorOn"`
}
