package conf

type Config struct {
	UserRPC    URPC       `mapstructure:"user_rpc" yaml:"user_rpc"`
	DB         DB         `mapstructure:"db" yaml:"db"`
	Log        Log        `mapstructure:"log" yaml:"log"`
	ConsulInfo ConsulInfo `mapstructure:"consul_info" yaml:"consul_info"`
	LocalInfo  LocalInfo  `mapstructure:"local_info" yaml:"local_info"`
}
