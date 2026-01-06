package global

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	"io"
	"option_service/conf"
)

var (
	DB          *gorm.DB
	RedisMutex  *redsync.Redsync
	Config      *conf.Config
	Tracer      opentracing.Tracer
	TracerClose io.Closer
)
