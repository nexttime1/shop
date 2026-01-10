package main

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"user_web/core"
	"user_web/flags"
	"user_web/global"
	"user_web/router"
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
	tracer, closer, err := core.InitTracer()
	if err != nil {
		zap.L().Error("启动jaeger失败")
		return
	}
	global.Tracer = tracer
	global.TracerClose = closer
	opentracing.SetGlobalTracer(global.Tracer)

	router.Router()
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
