package conf

type Config struct {
	System  System `mapstructure:"system"`
	UserRPC URPC   `mapstructure:"user_rpc"`
	Jwt     Jwt    `mapstructure:"jwt"`
	Redis   Redis  `mapstructure:"redis"`
	Code    Code   `mapstructure:"code"`
}
