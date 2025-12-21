package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"option_service/conf"
)

var (
	DB         *gorm.DB
	RedisMutex *redsync.Redsync
	Config     *conf.Config
)
