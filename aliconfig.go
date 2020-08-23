package main

import (
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

// ALiConfigClientConfig 阿里云 configClient需要的参数信息
type ALiConfigClientConfig struct {
	endpoint    string
	namespaceID string
	accessKey   string
	secretKey   string

	dataID string
	group  string
}

func (c *ALiConfigClientConfig) new() {

	endpoint := LocalViperConfig.GetString("aliyunConfig.endpoint")
	namespaceID := LocalViperConfig.GetString("aliyunConfig.namespaceid")
	accessKey := LocalViperConfig.GetString("aliyunConfig.accessKey")
	secretKey := LocalViperConfig.GetString("aliyunConfig.secretKey")

	dataID := LocalViperConfig.GetString("aliyunConfig.dataID")
	group := LocalViperConfig.GetString("aliyunConfig.group")

	fmt.Printf("endpoint: %s\nnamespaceID: %s\naccessKey: %s\nsecretKey: %s\ndataID: %s\ngroup: %s\n",
		endpoint,
		namespaceID,
		accessKey,
		secretKey,
		dataID,
		group,
	)

	c.endpoint = endpoint
	c.namespaceID = namespaceID
	c.accessKey = accessKey
	c.secretKey = secretKey
	c.dataID = dataID
	c.group = group
}

// 生成阿里云 configClient
func newConfigClient(c *LocalConfig) (config_client.IConfigClient, error) {
	clientConfig := constant.ClientConfig{
		Endpoint:       c.Endpoint + ":8080",
		NamespaceId:    c.NamespaceID,
		AccessKey:      c.AccessKey,
		SecretKey:      c.SecretKey,
		TimeoutMs:      3 * 1000,
		ListenInterval: 30 * 1000,
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})
	if err != nil {
		return nil, err
	}
	return configClient, nil
}

// func (c *ALiConfigClientConfig) watchConfigChange(ctx context.Context, ch chan int) error {
func (c *LocalConfig) watchConfigChange(ctx context.Context) error {
	aliConfigClient, err := newConfigClient(c)
	if err != nil {
		return err
	}

	fn := func(namespace, group, dataID, data string) {
		stringToViperConfig(data)
	}

	return aliConfigClient.ListenConfig(vo.ConfigParam{
		DataId:   c.DataID,
		Group:    c.Group,
		OnChange: fn,
	})
}
