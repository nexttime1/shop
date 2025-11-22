package core

import (
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"shop_api/global"
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
		fmt.Println("redis 连接失败")
		zap.S().Errorf("redis 连接失败  %s", err.Error())
	}
	zap.S().Info("redis 连接成功")
	return client
}
