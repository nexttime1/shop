package core

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"user_web/global"
)

func InitRedis() *redis.Client {
	r := global.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
	err := client.Ping().Err()
	if err != nil {
		zap.S().Errorf("redis 连接失败  %s", err.Error())
		return nil
	}
	zap.L().Info("redis 连接成功")
	return client
}
