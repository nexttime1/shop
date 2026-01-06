package conf

type Config struct {
	DB           DB         `mapstructure:"db" yaml:"db"`
	ConsulInfo   ConsulInfo `mapstructure:"consul_info" yaml:"consul_info"`
	LocalInfo    LocalInfo  `mapstructure:"local_info" yaml:"local_info"`
	NacosInfo    NacosInfo  `mapstructure:"nacos_info" yaml:"nacos_info"`
	Redis        Redis      `mapstructure:"redis" yaml:"redis"`
	GoodSrv      GoodSrv    `mapstructure:"good_srv" yaml:"good_srv"`
	InventorySrv GoodSrv    `mapstructure:"inventory_srv" yaml:"inventory_srv"`
	Alipay       Alipay     `mapstructure:"alipay" yaml:"alipay"`
	JaegerInfo   JaegerInfo `mapstructure:"jaeger_info" yaml:"jaeger_info"`
}
