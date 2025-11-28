package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"shop_api/conf"
	"shop_api/flags"
)

func ReadConf() *conf.Config {
	v := viper.New()
	v.SetConfigFile(flags.FileOption.File)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var c = new(conf.Config)
	err = v.Unmarshal(&c)

	if err != nil {
		panic(fmt.Sprintf("配置文件格式错误 ,%s", err))
	}

	// 读取 nacos 里的配置
	//  Nacos服务器配置
	serverConfigs := []constant.ServerConfig{
		{IpAddr: c.NacosInfo.Host, Port: c.NacosInfo.Port},
	}
	// 客户端配置
	clientConfig := constant.ClientConfig{
		Username:            c.NacosInfo.User,
		Password:            c.NacosInfo.Password,
		NamespaceId:         c.NacosInfo.Namespace, // id
		TimeoutMs:           5000,                  // 超时时间5秒
		NotLoadCacheAtStart: true,                  // 启动不加载缓存
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		zap.L().Error(err.Error())
		panic(err)
	}
	// 获得配置
	config, err := client.GetConfig(vo.ConfigParam{
		DataId: c.NacosInfo.DataId,
		Group:  c.NacosInfo.Group,
	})
	if err != nil {
		zap.L().Error(err.Error())
		panic(err)
	}
	err = yaml.Unmarshal([]byte(config), c)
	if err != nil {
		zap.L().Error(err.Error())
		panic(err)
	}

	fmt.Printf("读取配置文件 %s 成功\n", flags.FileOption.File)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", e.Name)
	})

	return c
}
