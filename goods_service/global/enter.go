package global

import (
	"goods_service/conf"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB

	Config *conf.Config
)
