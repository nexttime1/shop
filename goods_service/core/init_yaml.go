package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"goods_service/conf"
	"goods_service/flags"
	"goods_service/utils/nacos_get"
	"gopkg.in/yaml.v2"
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

	// 使用HTTP客户端直接调用Nacos API
	configContent, err := nacos_get.GetConfigWithTokenAuth(c.NacosInfo)
	if err != nil {
		panic("从Nacos获取配置失败: " + err.Error())
	}

	fmt.Printf("成功获取配置，内容长度: %d\n", len(configContent))
	fmt.Printf("配置内容预览:\n%s\n", configContent[:min(100, len(configContent))])

	// 解析配置
	err = yaml.Unmarshal([]byte(configContent), c)
	if err != nil {
		panic("解析Nacos配置失败: " + err.Error())
	}

	fmt.Printf("Nacos配置解析成功\n")

	// 配置文件监听
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", e.Name)
	})

	return c
}

// ReadConf111 第一版
func ReadConf111() *conf.Config {
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

	fmt.Printf("读取配置文件 %s 成功\n", flags.FileOption.File)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", e.Name)
	})

	return c
}
