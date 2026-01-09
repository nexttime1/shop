package conf

type Config struct {
	DB         DB         `mapstructure:"db" yaml:"db"`
	ConsulInfo ConsulInfo `mapstructure:"consul_info" yaml:"consul_info"`
	LocalInfo  LocalInfo  `mapstructure:"local_info" yaml:"local_info"`
	NacosInfo  NacosInfo  `mapstructure:"nacos_info" yaml:"nacos_info"`
	EsInfo     ConsulInfo `mapstructure:"es_info" yaml:"es_info"`
	JaegerInfo JaegerInfo `mapstructure:"jaeger_info" yaml:"jaeger_info"`
	Sentinel   Sentinel   `mapstructure:"sentinel" yaml:"sentinel"`
}
