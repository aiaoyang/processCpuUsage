package configs

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/golang/glog"
	"github.com/nacos-group/nacos-sdk-go/clients"
	nacos "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// AliyunAcm 访问 ·阿里云配置中心· 认证配置
type AliyunAcm struct {
	Endpoint    string `yaml:"endpoint"`
	NamespaceID string `yaml:"namespaceID"`
	AccessKey   string `yaml:"accessKey"`
	SecretKey   string `yaml:"secretKey"`

	DataID string `yaml:"dataID"`
	Group  string `yaml:"group"`
}

// InfluxDBConfig 访问 ·influxdb· 认证配置
type InfluxDB struct {
	DBName    string `yaml:"dbName"`
	TableName string `yaml:"tableName"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
}

// Config 程序启动配置
type Config struct {
	AliyunAcm `yaml:"aliyunAcm"`
	InfluxDB  `yaml:"influxDB"`

	// 部署环境 ·dev· or ·release·
	Env string `yaml:"env"`
}

func init() {

}

// 连接阿里云配置中心
func newConfigClient(c *Config) (acmClient nacos.IConfigClient) {
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

	acmClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})
	// acmClient.PublishConfig()
	if err != nil {
		glog.Fatal(errors.Wrap(err, "connect to aliyun acm failed"))
	}

	return
}
func LoadConfig(cfgPath string) (c Config) {
	c = Config{}
	f, err := os.Open(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func LoadAcmConfig(ctx context.Context, c Config, acmCh chan *AcmInnerConfig) (err error) {

	acmClient := newConfigClient(&c)

	go func() { // 获取第一次配置
		conf, err := acmClient.GetConfig(vo.ConfigParam{
			DataId:   c.DataID,
			Group:    c.Group,
			OnChange: nil,
		})
		if err != nil {
			acmCh <- nil
			return
		}

		cfg, err := stringToViperConfig(conf)
		if err != nil {
			acmCh <- nil
			return
		}

		acmCh <- &cfg
	}()

	go func() {
		err = c.watchAcm(ctx, acmClient, acmCh)
		if err != nil {
			acmCh <- nil
			return
		}
	}()

	return
}

func (c *Config) watchAcm(ctx context.Context, acmClient nacos.IConfigClient, configCh chan *AcmInnerConfig) (err error) {

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
}
