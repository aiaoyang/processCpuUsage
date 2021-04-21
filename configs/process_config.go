package configs

// Config 进程状态配置
type ProcessStat struct {
	ProcessInfo `yaml:"processInfo"`
}

// ProcessInfo 监控进程的信息
type ProcessInfo struct {
	HostName string   `yaml:"hostName"`
	PIDS     []int    `yaml:"pids"`
	Names    []string `yaml:"names"`
	Status   int      `yaml:"status"`
}
