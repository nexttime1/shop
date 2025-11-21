package core

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"shop_api/conf"
	"shop_api/flags"
)

func ReadConf() *conf.Config {
	file, err := os.ReadFile(flags.FileOption.File)
	if err != nil {
		panic(err)
	}
	var c = new(conf.Config)
	err = yaml.Unmarshal(file, c)
	if err != nil {
		panic(fmt.Sprintf("yaml配置文件格式错误 ,%s", err))
	}

	fmt.Printf("读取配置文件 %s 成功\n", flags.FileOption.File)

	return c
}
