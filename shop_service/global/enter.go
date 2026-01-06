package global

import (
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	"io"
	"shop_service/conf"
)

var (
	DB          *gorm.DB
	Tracer      opentracing.Tracer
	TracerClose io.Closer
	Config      *conf.Config
)
