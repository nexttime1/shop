package conf

type Config struct {
	UserRPC URPC `yaml:"user_rpc"`
	DB      DB   `yaml:"db"`
	Log     Log  `yaml:"log"`
}
