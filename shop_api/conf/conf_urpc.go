package conf

import "fmt"

type URPC struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

func (s URPC) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
