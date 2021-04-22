package configs

import (
	"bytes"
	"context"

	"github.com/golang/glog"
	"github.com/nacos-group/nacos-sdk-go/clients"
	nacos "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// InfluxDBConfig 本地配置中influxdb相关配置
type InfluxDBConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var (
	// MyLocalConfig config 本地配置
	GlobalConfig = Config{}
)

// Config 程序启动配置
type Config struct {
	AliyunAcm      `yaml:"aliyunAcm"`
	InfluxDBConfig `yaml:"influxDB"`
	Env            string `yaml:"env"`
}

func init() {

	// 阿里云配置文件初始化
	// AliyunConfigSrv = viper.New()
	// AliyunConfigSrv.SetConfigType("yaml")

	// // 本地配置文件初始化
	// LocalViperConfig = viper.New()
	// LocalViperConfig.SetConfigType("yaml")
	// LocalViperConfig.SetConfigFile("./config.yaml")

	// err := LocalViperConfig.ReadInConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = LocalViperConfig.Unmarshal(&MyLocalConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ctx := context.TODO()
	// MyLocalConfig.watchConfigChange(ctx)

	// fmt.Printf("local Config : %v\n", MyLocalConfig)
}

// 连接阿里云配置中心
func newConfigClient(c *Config) (acmClient nacos.IConfigClient, err error) {
	clientConfig := constant.ClientConfig{
		Endpoint:        c.Endpoint + ":8080",
		NamespaceId:     c.NamespaceID,
		AccessKey:       c.AccessKey,
		SecretKey:       c.SecretKey,
		TimeoutMs:       3 * 1000,
		ListenInterval:  30 * 1000,
		UpdateThreadNum: 1,
		LogLevel:        "warn",
	}

	acmClient, err = clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})
	if err != nil {
		glog.Fatal(errors.Wrap(err, "connect to aliyun acm failed"))
	}

	return
}

func LoadConfig(ctx context.Context, c Config) (configCh chan *AcmInnerConfig, err error) {

	acmClient, _ := newConfigClient(&c)

	go func() {
		err = c.watchAcmChanged(ctx, acmClient, configCh)
		if err != nil {
			configCh <- nil
			return
		}
	}()

	go func() { // 获取第一次配置
		conf, err := acmClient.GetConfig(vo.ConfigParam{
			DataId:   c.DataID,
			Group:    c.Group,
			OnChange: nil,
		})
		if err != nil {
			configCh <- nil
			return
		}

		cfg, err := stringToViperConfig(conf)
		if err != nil {
			configCh <- nil
			return
		}

		configCh <- &cfg
	}()

	return

}

func (c *Config) watchAcmChanged(ctx context.Context, acmClient nacos.IConfigClient, configCh chan *AcmInnerConfig) (err error) {

	onChangedAction := func(namespace, group, dataID, data string) {
		acmInnerCfg, err := stringToViperConfig(data)
		if err != nil {
			glog.Errorf("config changed err: %s\n", err)
			return
		}
		configCh <- &acmInnerCfg
	}

	// 开始监听配置变化
	return acmClient.ListenConfig(vo.ConfigParam{
		DataId:   c.DataID,
		Group:    c.Group,
		OnChange: onChangedAction,
	})

}

// 从阿里云获得的配置文件内容转换为viper的配置文件内容
func stringToViperConfig(s string) (acmCfg AcmInnerConfig, err error) {
	cfg := AcmInnerConfig{}
	err = yaml.NewDecoder(bytes.NewBuffer([]byte(s))).Decode(&cfg)
	if err != nil {
		return
	}
	return
	// err := AliyunConfigSrv.ReadConfig(bytes.NewBuffer([]byte(s)))
	// if err != nil {
	// 	// log.Fatalf("err: %v", err)
	// 	panic(err)
	// }
	// // 配置中心配置改变，同步到进程监控配置
	// err = AliyunConfigSrv.Unmarshal(&NeedMonitorProcessInfo)
	// if err != nil {
	// 	panic(err)
	// }
	// NeedMonitorProcessInfo.PIDS = sysusage.GetProcessPID(NeedMonitorProcessInfo.Names...)
	// log.Printf("config : %v\n", NeedMonitorProcessInfo)
	// log.Printf("pid change to : %d\n", NeedMonitorProcessInfo.PIDS)
}
