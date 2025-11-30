package conf

import "fmt"

type ConsulInfo struct {
	Port int    `mapstructure:"port" yaml:"port"`
	Addr string `mapstructure:"addr" yaml:"addr"`
	Name string `mapstructure:"name" yaml:"name"`
}

func (c ConsulInfo) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}
