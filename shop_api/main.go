package main

import (
	"shop_api/core"
	"shop_api/flags"
	"shop_api/global"
	"shop_api/router"
)

func main() {
	flags.Parse()
	core.InitZap()
	global.Config = core.ReadConf()
	//fmt.Println(global.Config)
	global.Redis = core.InitRedis()
	router.Router()

}
