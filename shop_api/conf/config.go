package conf

type Config struct {
	System  System `mapstructure:"system"`
	UserRPC URPC   `mapstructure:"user_rpc"`
}
