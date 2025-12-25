package global

import (
	"github.com/olivere/elastic/v7"
	"goods_service/conf"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	EsClient *elastic.Client
	Config   *conf.Config
)
