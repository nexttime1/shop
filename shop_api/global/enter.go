package global

import (
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"io"
	"shop_api/conf"
)

var (
	Config      *conf.Config
	Redis       *redis.Client
	LevelFlag   bool
	Tracer      opentracing.Tracer
	TracerClose io.Closer
)
