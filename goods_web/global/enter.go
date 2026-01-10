package global

import (
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"goods_web/conf"
	"io"
)

var (
	Config      *conf.Config
	Redis       *redis.Client
	LevelFlag   bool
	Tracer      opentracing.Tracer
	TracerClose io.Closer
)
