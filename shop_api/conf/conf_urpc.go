package conf

import "fmt"

type URPC struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

func (s URPC) GetAddr() string {
	return fmt.Sprintf("%s:%d", s.IP, s.Port)
}
