package conf

import "fmt"

type System struct {
	IP      string `mapstructure:"ip"`
	Port    int    `mapstructure:"port"`
	GinMode string `mapstructure:"gin_mode"`
	Env     string `mapstructure:"env"`
}

func (s System) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
