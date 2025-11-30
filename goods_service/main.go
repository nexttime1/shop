package main

import (
	"fmt"

	"goods_service/core"
	"goods_service/flags"
	"goods_service/global"
)

func main() {
	flags.Parse() //解析 yaml文件
	core.InitZap()
	global.Config = core.ReadConf()
	//fmt.Println(global.Config)
	global.DB = core.InitDB()

	//core.InitLogrus()
	flags.Run()
	//err := core.InitRPC()

	//if err != nil {
	//	return
	//}
	fmt.Println("运行成功")

}
