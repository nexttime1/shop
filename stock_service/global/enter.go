package global

import (
	"gorm.io/gorm"
	"stock_service/conf"
)

var (
	DB *gorm.DB

	Config *conf.Config
)
