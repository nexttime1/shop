package global

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/go-redsync/redsync/v4"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	"io"
	"order_service/conf"
)

var (
	DB          *gorm.DB
	RedisMutex  *redsync.Redsync
	Config      *conf.Config
	Producer    rocketmq.TransactionProducer
	Tracer      opentracing.Tracer
	TracerClose io.Closer
)
