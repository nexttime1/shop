package main

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"time"
)

func main() {
	// 1. 配置Jaeger
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, // 常量采样，1表示全量采样
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,                   // 打印span日志
			LocalAgentHostPort:  "192.168.163.132:6831", // 本地Jaeger Agent地址（关键修改）
			BufferFlushInterval: 1 * time.Second,        // 主动刷新上报间隔
		},
		ServiceName: "ceshi", // 服务名
	}

	// 2. 创建tracer和closer
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger), // 日志输出
	)
	if err != nil {
		panic("创建tracer失败: " + err.Error())
	}
	// 关键修改：创建后立即defer关闭，确保最后执行（必须在err判断后）
	defer func() {
		closer.Close()
		time.Sleep(2 * time.Second) // 给上报留时间
	}()

	// 3. 设置全局tracer
	opentracing.SetGlobalTracer(tracer)

	// 4. 创建并结束span
	span := tracer.StartSpan("main")
	defer span.Finish() // 关键修改：defer结束span，避免遗漏

	// 模拟业务逻辑
	time.Sleep(10 * time.Second)

	// 等待上报完成（可选，确保span数据发送）
	time.Sleep(1 * time.Second)
}
