package main

import (
	"fmt"
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
	fmt.Println("运行成功")
}
