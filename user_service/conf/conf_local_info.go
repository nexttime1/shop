package conf

import "fmt"

type LocalInfo struct {
	Port int    `mapstructure:"port" yaml:"port"`
	Addr string `mapstructure:"addr" yaml:"addr"`
}

func (l LocalInfo) GetAddr() string {
	return fmt.Sprintf("%s:%d", l.Addr, l.Port)
}
