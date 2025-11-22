package global

import (
	"github.com/go-redis/redis"
	"shop_api/conf"
)

var (
	Config    *conf.Config
	Redis     *redis.Client
	LevelFlag bool
)
