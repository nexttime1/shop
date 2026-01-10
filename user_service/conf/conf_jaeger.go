package conf

import "fmt"

type JaegerInfo struct {
	Host        string `mapstructure:"host" yaml:"host"`
	Port        uint64 `mapstructure:"port" yaml:"port"`
	ServiceName string `mapstructure:"service_name" yaml:"service_name"`
}

func (info JaegerInfo) Addr() string {
	return fmt.Sprintf("%s:%d", info.Host, info.Port)
}
