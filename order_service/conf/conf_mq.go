package conf

import (
	"fmt"
)

type RocketMQ struct {
	Host              string `mapstructure:"host" yaml:"host"`
	Port              uint64 `mapstructure:"port" yaml:"port"`
	GroupName         string `mapstructure:"group_name" yaml:"group_name"`
	Topic             string `mapstructure:"topic" yaml:"topic"`
	ConsumerGroupName string `mapstructure:"consumer_group_name" yaml:"consumer_group_name"`
	ConsumerSubscribe string `mapstructure:"consumer_subscribe" yaml:"consumer_subscribe"`
	ConsumerTopic     string `mapstructure:"consumer_topic" yaml:"consumer_topic"`
	MaxRetryTimes     int    `mapstructure:"max_retry_times" yaml:"max_retry_times"`
	BaseRetryDelay    int    `mapstructure:"base_retry_delay" yaml:"base_retry_delay"`
}

func (info RocketMQ) Addr() string {
	return fmt.Sprintf("%s:%d", info.Host, info.Port)
}
