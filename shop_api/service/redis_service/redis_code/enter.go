package redis_code

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"shop_api/global"
)

type RedisCode struct {
	Index string
}

func NewRegisterIndex() *RedisCode {
	return &RedisCode{
		Index: "register_code",
	}
}

func NewRecallIndex() *RedisCode {
	return &RedisCode{
		Index: "recall_code",
	}
}

func (r RedisCode) StorageCode(code, phone string) error {
	err := global.Redis.HSet(r.Index, phone, code).Err()
	if err != nil {
		return fmt.Errorf("failed to store verify code: %w", err)
	}

	return nil
}

func (r RedisCode) GetCode(phone string) (string, error) {
	code, err := global.Redis.HGet(r.Index, phone).Result()
	if err != nil {
		zap.S().Errorf(err.Error())
		return "", errors.New("不存在验证码")
	}
	return code, nil
}

func (r RedisCode) DeleteCode(phone string) error {
	_, err := global.Redis.HDel(r.Index, phone).Result()
	if err != nil {
		zap.S().Errorf(err.Error())
		return fmt.Errorf("删除使用的验证码失败")
	}
	return nil
}

func (r RedisCode) Clear() {
	global.Redis.Del(r.Index)

}
