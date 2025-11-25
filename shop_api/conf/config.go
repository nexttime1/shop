package conf

type Config struct {
	System     System     `mapstructure:"system"`
	Jwt        Jwt        `mapstructure:"jwt"`
	Redis      Redis      `mapstructure:"redis"`
	Code       Code       `mapstructure:"code"`
	ConsulInfo ConsulInfo `mapstructure:"consul_info"`
}
