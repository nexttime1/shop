package main

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"user_service/core"
	"user_service/flags"
	"user_service/global"
)

func main() {
	flags.Parse() //解析 yaml文件
	core.InitZap()
	global.Config = core.ReadConf()

	global.DB = core.InitDB()

	flags.Run()
	tracer, closer, err := core.InitTracer()
	if err != nil {
		zap.L().Error("tracer 初始化失败", zap.Error(err))
		panic(err)
	}
	global.Tracer = tracer
	global.TracerClose = closer
	opentracing.SetGlobalTracer(tracer)

	client := core.NewConsulRegister()
	err = client.Register()
	if err != nil {
		zap.L().Error("注册失败", zap.Error(err))
		panic(err)
	}
	// ctrl + C 自动注销 刚注册的consul  监听
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit // 阻塞
	err = client.Deregister()
	global.TracerClose.Close()
	if err != nil {
		zap.L().Error("服务注销失败", zap.Error(err))
		return
	}

	zap.S().Info("服务注销成功")

}
