package conf

type Config struct {
	UserRPC    URPC       `mapstructure:"user_rpc"`
	DB         DB         `mapstructure:"db"`
	Log        Log        `mapstructure:"log"`
	ConsulInfo ConsulInfo `mapstructure:"consul_info"`
	LocalInfo  LocalInfo  `mapstructure:"local_info"`
}
