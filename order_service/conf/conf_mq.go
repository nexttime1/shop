package conf

import (
	"fmt"
)

type RocketMQ struct {
	Host                     string `mapstructure:"host" yaml:"host"`
	Port                     uint64 `mapstructure:"port" yaml:"port"`
	TransGroupName           string `mapstructure:"trans_group_name" yaml:"trans_group_name"`
	TransactionTopic         string `mapstructure:"transaction_topic" yaml:"transaction_topic"`
	DelayGroupName           string `mapstructure:"delay_group_name" yaml:"delay_group_name"`
	DelayTopic               string `mapstructure:"delay_topic" yaml:"delay_topic"`
	TimeOutConsumerGroupName string `mapstructure:"timeout_consumer_group_name" yaml:"timeout_consumer_group_name"`
	TimeOutTopic             string `mapstructure:"timeout_topic" yaml:"timeout_topic"`
	StockTimeoutTopic        string `mapstructure:"stock_timeout_topic" yaml:"stock_timeout_topic"`
	MaxRetryTimes            int32  `mapstructure:"max_retry_times" yaml:"max_retry_times"`
	BaseRetryDelay           int    `mapstructure:"base_retry_delay" yaml:"base_retry_delay"`
}

func (info RocketMQ) Addr() string {
	return fmt.Sprintf("%s:%d", info.Host, info.Port)
}
