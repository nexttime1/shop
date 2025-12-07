package global

import (
	"github.com/go-redis/redis"
	"oss_web/conf"
)

var (
	Config    *conf.Config
	Redis     *redis.Client
	LevelFlag bool
)
