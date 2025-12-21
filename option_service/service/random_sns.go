package service

import (
	"fmt"
	"math/rand"
	"time"
)

func RandomSns(userID int32) string {
	now := time.Now()
	rand.Seed(now.UnixNano()) //毫秒级go
	id := rand.Intn(90) + 10  // 两位随机数
	OrderSns := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userID, id)

	return OrderSns
}
