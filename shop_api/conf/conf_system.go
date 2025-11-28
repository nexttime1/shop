package conf

import "fmt"

type System struct {
	IP      string `mapstructure:"ip" yaml:"ip"`
	Port    int    `mapstructure:"port" yaml:"port"`
	GinMode string `mapstructure:"gin_mode" yaml:"gin_mode"`
	Env     string `mapstructure:"env" yaml:"env"`
}

func (s System) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
