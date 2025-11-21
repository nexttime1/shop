package conf

type Config struct {
	System  System `yaml:"system"`
	UserRPC URPC   `yaml:"user_rpc"`
}
