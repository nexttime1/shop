package conf

import "fmt"

type ConsulInfo struct {
	Port int      `mapstructure:"port" yaml:"port"`
	Addr string   `mapstructure:"addr" yaml:"addr"`
	Name string   `mapstructure:"name" yaml:"name"`
	Tags []string `mapstructure:"tags" yaml:"tags"`
}

func (c ConsulInfo) GetAddr() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}
