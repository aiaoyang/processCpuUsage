package configs

import (
	"context"
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients"
	nacos "github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// LocalConfig 本地配置文件
type LocalConfig struct {
	AliYunConfig   `yaml:"aliyunConfig"`
	InfluxDBConfig `yaml:"influxDBConfig"`
	Env            string `yaml:"env"`
}

// 连接阿里云配置中心
func newConfigClient(c *LocalConfig) (nacos.IConfigClient, error) {
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

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})
	if err != nil {
		return nil, err
	}
	return configClient, nil
}

func (c *LocalConfig) watchConfigChange(ctx context.Context) error {
	aliConfigClient, err := newConfigClient(c)
	if err != nil {
		return err
	}

	fn := func(namespace, group, dataID, data string) {
		stringToViperConfig(data)
	}

	// 获取第一次配置
	conf, err := aliConfigClient.GetConfig(vo.ConfigParam{
		DataId:   c.DataID,
		Group:    c.Group,
		OnChange: fn,
	})
	if err != nil {
		log.Fatal(err)
	}
	stringToViperConfig(conf)

	// 开始监听配置变化
	return aliConfigClient.ListenConfig(vo.ConfigParam{
		DataId:   c.DataID,
		Group:    c.Group,
		OnChange: fn,
	})
}
