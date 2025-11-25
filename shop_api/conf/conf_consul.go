package conf

import "fmt"

type ConsulInfo struct {
	Port int    `mapstructure:"port"`
	Addr string `mapstructure:"addr"`
	Name string `mapstructure:"name"`
}

func (c ConsulInfo) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}
