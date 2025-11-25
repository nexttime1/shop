package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"shop_service/core"
	"shop_service/flags"
	"shop_service/global"
)

func main() {
	flags.Parse() //解析 yaml文件
	global.Config = core.ReadConf()
	global.DB = core.InitDB()
	core.InitLogrus()
	flags.Run()
	err := core.InitRPC()

	if err != nil {
		logrus.Errorf("init rpc error: %v", err)
		fmt.Println(err)
	}
	fmt.Println("运行成功")

}
