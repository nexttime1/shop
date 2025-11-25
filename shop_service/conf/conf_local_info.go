package conf

import "fmt"

type LocalInfo struct {
	Port int    `mapstructure:"port"`
	Addr string `mapstructure:"addr"`
}

func (l LocalInfo) GetAddr() string {
	return fmt.Sprintf("%s:%d", l.Addr, l.Port)
}
