package conf

type Config struct {
	DB         DB         `mapstructure:"db" yaml:"db"`
	ConsulInfo ConsulInfo `mapstructure:"consul_info" yaml:"consul_info"`
	LocalInfo  LocalInfo  `mapstructure:"local_info" yaml:"local_info"`
	NacosInfo  NacosInfo  `mapstructure:"nacos_info" yaml:"nacos_info"`
	Redis      Redis      `mapstructure:"redis" yaml:"redis"`
}
