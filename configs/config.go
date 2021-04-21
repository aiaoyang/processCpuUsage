package configs

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aiaoyang/processCpuUsage/pkg/sysusage"
	"github.com/spf13/viper"
)

// InfluxDBConfig 本地配置中influxdb相关配置
type InfluxDBConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	// NeedMonitorProcessInfo 被监控的进程状态信息
	NeedMonitorProcessInfo = ProcessStat{ProcessInfo{Status: 0}}

	// AliyunConfigSrv 阿里云配置文件
	AliyunConfigSrv *viper.Viper

	// LocalViperConfig viper 本地配置
	LocalViperConfig *viper.Viper

	// MyLocalConfig config 本地配置
	MyLocalConfig = LocalConfig{}
)

func init() {
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

	ctx := context.TODO()
	MyLocalConfig.watchConfigChange(ctx)

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
	NeedMonitorProcessInfo.PIDS = sysusage.GetProcessPID(NeedMonitorProcessInfo.Names...)
	log.Printf("config : %v\n", NeedMonitorProcessInfo)
	log.Printf("pid change to : %d\n", NeedMonitorProcessInfo.PIDS)
}
