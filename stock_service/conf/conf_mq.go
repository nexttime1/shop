package conf

import (
	"fmt"
)

type RocketMQ struct {
	Host              string `mapstructure:"host" yaml:"host"`
	Port              uint64 `mapstructure:"port" yaml:"port"`
	ConsumerGroupName string `mapstructure:"consumer_group_name" yaml:"consumer_group_name"`
	TransactionTopic  string `mapstructure:"transaction_topic" yaml:"transaction_topic"`
	StockTimeoutTopic string `mapstructure:"stock_timeout_topic" yaml:"stock_timeout_topic"`
	MaxRetryTimes     int32  `mapstructure:"max_retry_times" yaml:"max_retry_times"`
	BaseRetryDelay    int    `mapstructure:"base_retry_delay" yaml:"base_retry_delay"`
}

func (info RocketMQ) Addr() string {
	return fmt.Sprintf("%s:%d", info.Host, info.Port)
}
