package main

import (
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

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

	flags.Run()
	client := core.NewConsulRegister()

	err := client.Register()
	if err != nil {
		zap.L().Error("注册失败", zap.Error(err))
		panic(err)
	}
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
