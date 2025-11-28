package conf

type Config struct {
	System     System     `mapstructure:"system" yaml:"system"`
	Jwt        Jwt        `mapstructure:"jwt" yaml:"jwt"`
	Redis      Redis      `mapstructure:"redis" yaml:"redis"`
	Code       Code       `mapstructure:"code" yaml:"code"`
	ConsulInfo ConsulInfo `mapstructure:"consul_info" yaml:"consul_info"`
	NacosInfo  NacosInfo  `mapstructure:"nacos_info" yaml:"nacos_info"`
}
