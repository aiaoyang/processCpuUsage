package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/aiaoyang/processCpuUsage/sysusage"
	"github.com/spf13/viper"
)

// Config 进程状态配置
type Config struct {
	ProcessInfo `yaml:"processInfo"`
}

// ProcessInfo 监控进程的信息
type ProcessInfo struct {
	HostName string `yaml:"hostName"`
	PIDS     []int  `yaml:"pid"`
	Name     string `yaml:"name"`
	Status   int    `yaml:"status"`
}

// LocalConfig 本地配置文件
type LocalConfig struct {
	AliYunConfig   `yaml:"aliyunConfig"`
	InfluxDBConfig `yaml:"influxDBConfig"`
	Env            string `yaml:"env"`
}

// AliYunConfig 本地配置中阿里云相关访问权限字段设置
type AliYunConfig struct {
	Endpoint    string `yaml:"endpoint"`
	NamespaceID string `yaml:"namespaceID"`
	AccessKey   string `yaml:"accessKey"`
	SecretKey   string `yaml:"secretKey"`

	DataID string `yaml:"dataID"`
	Group  string `yaml:"group"`
}

// InfluxDBConfig 本地配置中influxdb相关配置
type InfluxDBConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	// NeedMonitorProcessInfo 被监控的进程状态信息
	NeedMonitorProcessInfo = Config{ProcessInfo{Status: 0}}
)
var (
	// AliyunConfigSrv 阿里云配置文件
	AliyunConfigSrv *viper.Viper
)

// LocalViperConfig 本地配置文件
var (
	// viper 本地配置
	LocalViperConfig *viper.Viper

	// config 本地配置
	MyLocalConfig = LocalConfig{}
)

func initViper() {
	// 阿里云配置文件初始化
	AliyunConfigSrv = viper.New()
	AliyunConfigSrv.SetConfigType("yaml")

	// 本地配置文件初始化
	LocalViperConfig = viper.New()
	LocalViperConfig.SetConfigType("yaml")
	LocalViperConfig.SetConfigFile("./config.yaml")

	err := LocalViperConfig.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = LocalViperConfig.Unmarshal(&MyLocalConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("local Config : %v\n", MyLocalConfig)
}

// 从阿里云获得的配置文件内容转换为viper的配置文件内容
func stringToViperConfig(s string) {
	err := AliyunConfigSrv.ReadConfig(bytes.NewBuffer([]byte(s)))
	if err != nil {
		// log.Fatalf("err: %v", err)
		panic(err)
	}
	// 配置中心配置改变，同步到进程监控配置
	err = AliyunConfigSrv.Unmarshal(&NeedMonitorProcessInfo)
	if err != nil {
		panic(err)
	}
	NeedMonitorProcessInfo.PIDS = sysusage.GetProcessPID(NeedMonitorProcessInfo.Name)
	log.Printf("pid change to : %d\n", NeedMonitorProcessInfo.PIDS)
}
