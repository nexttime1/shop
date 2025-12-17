package core

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"stock_service/global"
	"time"
)

func InitRedisMutex() *redsync.Redsync {
	// Create a pool with go-redis (or redigo) which is the pool redsync will
	// use while communicating with Redis. This can also be any pool that
	// implements the `redis.Pool` interface.
	Addr := global.Config.Redis.Addr
	client := goredislib.NewClient(&goredislib.Options{
		Addr:         Addr,
		Password:     "",               // 无密码
		DB:           0,                // 默认库
		PoolSize:     20,               // 连接池大小（根据并发调整）
		MinIdleConns: 5,                // 最小空闲连接，避免频繁创建连接
		DialTimeout:  5 * time.Second,  // 连接超时
		ReadTimeout:  3 * time.Second,  // 读超时
		WriteTimeout: 3 * time.Second,  // 写超时
		PoolTimeout:  30 * time.Second, // 空闲连接超时（小于Redis的tcp-keepalive）
	})

	fmt.Println(Addr)
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redsync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)
	zap.S().Info("redsync init redis 成功")

	return rs
}
