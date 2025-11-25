package conf

import "fmt"

type DB struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       string `mapstructure:"db"`
	Debug    bool   `mapstructure:"debug"`  //打印全部日志
	Source   string `mapstructure:"source"` //数据库类型 pgsql mysql
}

func (d DB) DSN() string {
	timeout := "10s"
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", d.User, d.Password, d.Host, d.Port, d.DB, timeout)
}

func (d DB) Empty() bool {
	return d.User == "" && d.Password == "" && d.Host == "" && d.Port == 0
}
