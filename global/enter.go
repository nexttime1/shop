package global

import (
	"gorm.io/gorm"
	"shop_service/conf"
)

var (
	DB     *gorm.DB
	Config *conf.Config
)
