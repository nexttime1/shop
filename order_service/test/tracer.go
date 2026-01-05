package main

import (
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func main() {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, // 所有的都要 跟踪
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "192.168.163.132:6831",
		},
		ServiceName: "shop_test",
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	parentScan := tracer.StartSpan("main")
	span := tracer.StartSpan("funcA", opentracing.ChildOf(parentScan.Context())) // 绑定父子
	time.Sleep(1 * time.Second)
	span.Finish()

	span2 := tracer.StartSpan("funcB", opentracing.ChildOf(parentScan.Context())) // 绑定父子
	time.Sleep(500 * time.Millisecond)
	span2.Finish()

}
