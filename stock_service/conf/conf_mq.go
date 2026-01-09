package conf

import (
	"fmt"
)

type RocketMQ struct {
	Host              string `mapstructure:"host" yaml:"host"`
	Port              uint64 `mapstructure:"port" yaml:"port"`
	ConsumerGroupName string `mapstructure:"consumer_group_name" yaml:"consumer_group_name"`
	ConsumerSubscribe string `mapstructure:"consumer_subscribe" yaml:"consumer_subscribe"`
}

func (info RocketMQ) Addr() string {
	return fmt.Sprintf("%s:%d", info.Host, info.Port)
}
