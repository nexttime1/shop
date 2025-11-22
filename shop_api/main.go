package main

import (
	"shop_api/core"
	"shop_api/flags"
	"shop_api/global"
	"shop_api/router"
)

func main() {
	flags.Parse()
	global.Config = core.ReadConf()
	global.Redis = core.InitRedis()
	core.InitZap()

	router.Router()

}
