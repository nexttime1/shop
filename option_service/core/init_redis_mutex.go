package core

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"option_service/global"
)

func InitRedisMutex() *redsync.Redsync {
	// Create a pool with go-redis (or redigo) which is the pool redsync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	ip := global.Config.LocalInfo.Addr
	port := global.Config.Redis.Addr
	Addr := ip + ":" + port
	client := goredislib.NewClient(&goredislib.Options{
		Addr: Addr,
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redsync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)
	zap.S().Info("redsync init redis 成功")

	return rs
}
