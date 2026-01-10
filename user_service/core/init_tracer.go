package core

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"go.uber.org/zap"

	"io"
	"user_service/global"
)

func InitTracer() (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, // 所有的都要 跟踪
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: global.Config.JaegerInfo.Addr(),
		},
		ServiceName: global.Config.JaegerInfo.ServiceName,
	}
	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaegerlog.StdLogger))
	if err != nil {
		zap.S().Error(err)
	}

	return tracer, closer, err

}
