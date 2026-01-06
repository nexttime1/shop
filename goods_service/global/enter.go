package global

import (
	"github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
	"goods_service/conf"
	"gorm.io/gorm"
	"io"
)

var (
	DB          *gorm.DB
	EsClient    *elastic.Client
	Config      *conf.Config
	Tracer      opentracing.Tracer
	TracerClose io.Closer
)
