package conf

import "fmt"

type RocketMQ struct {
	Host      string `mapstructure:"host" yaml:"host"`
	Port      uint64 `mapstructure:"port" yaml:"port"`
	GroupName string `mapstructure:"group_name" yaml:"group_name"`
}

func (info RocketMQ) Addr() string {
	return fmt.Sprintf("%s:%d", info.Host, info.Port)
}
