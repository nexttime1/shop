package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

	fmt.Printf("读取配置文件 %s 成功\n", flags.FileOption.File)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", e.Name)
	})

	return c
}
