package conf

type Config struct {
	System       System     `mapstructure:"system" yaml:"system"`
	Jwt          Jwt        `mapstructure:"jwt" yaml:"jwt"`
	GoodSrv      GoodSrv    `mapstructure:"good_srv" yaml:"good_srv"`
	InventorySrv GoodSrv    `mapstructure:"inventory_srv" yaml:"inventory_srv"`
	Redis        Redis      `mapstructure:"redis" yaml:"redis"`
	ConsulInfo   ConsulInfo `mapstructure:"consul_info" yaml:"consul_info"`
	NacosInfo    NacosInfo  `mapstructure:"nacos_info" yaml:"nacos_info"`
	Alipay       Alipay     `mapstructure:"alipay" yaml:"alipay"`
	JaegerInfo   JaegerInfo `mapstructure:"jaeger_info" yaml:"jaeger_info"`
	Sentinel     Sentinel   `mapstructure:"sentinel" yaml:"sentinel"`
}
