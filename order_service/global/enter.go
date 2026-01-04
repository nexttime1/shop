package global

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"order_service/conf"
)

var (
	DB         *gorm.DB
	RedisMutex *redsync.Redsync
	Config     *conf.Config
	Producer   rocketmq.TransactionProducer
)
