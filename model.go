package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config process config
type Config struct {
	ProcessInfo `yaml:"processInfo"`
}

// ProcessInfo 监控进程的信息
type ProcessInfo struct {
	HostName string   `yaml:"hostName"`
	PIDS     []string `yaml:"pid"`
	Name     string   `yaml:"name"`
	Status   int      `yaml:"status"`
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
	DataID      string `yaml:"dataID"`
	Group       string `yaml:"group"`
}

// InfluxDBConfig 本地配置中influxdb相关配置
type InfluxDBConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	// VCloud 阿里云配置文件
	VCloud *viper.Viper
	// GConfig global config setting
	GConfig = Config{ProcessInfo{Status: 0}}
)

// VLocal 本地配置文件
var (
	VLocal *viper.Viper
	CLocal = LocalConfig{}
)

func initViper() {
	// 阿里云配置文件初始化
	VCloud = viper.New()
	VCloud.SetConfigType("yaml")

	// 本地配置文件初始化
	VLocal = viper.New()
	VLocal.SetConfigType("yaml")
	VLocal.SetConfigFile("./config.yaml")
	err := VLocal.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = VLocal.Unmarshal(&CLocal)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("local Config : %v\n", CLocal)
}

// 从阿里云获得的配置文件内容转换为viper的配置文件内容
func stringToViperConfig(v *viper.Viper, s string) {
	err := v.ReadConfig(bytes.NewBuffer([]byte(s)))
	if err != nil {
		// log.Fatalf("err: %v", err)
		panic(err)
	}
	err = v.Unmarshal(&GConfig)
	if err != nil {
		panic(err)
	}
	GConfig.PIDS = getProcessPID(GConfig.Name)
}
