package main

import (
	"go.uber.org/zap"
	"goods_api/core"
	"goods_api/flags"
	"goods_api/global"
	"goods_api/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flags.Parse()
	core.InitZap()
	global.Config = core.ReadConf()
	//fmt.Println(global.Config)
	global.Redis = core.InitRedis()
	client := core.NewConsulRegister()
	err := client.Register()
	if err != nil {
		zap.L().Error("注册失败", zap.Error(err))
		panic(err)
	}
	router.Router()
	// ctrl + C 自动注销 刚注册的consul  监听
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit // 阻塞
	err = client.Deregister()
	if err != nil {
		zap.L().Error("服务注销失败", zap.Error(err))
		return
	}

	zap.S().Info("服务注销成功")

}
