package conf

type NacosInfo struct {
	Host      string `mapstructure:"host" yaml:"host"`
	Port      uint64 `mapstructure:"port" yaml:"port"`
	User      string `mapstructure:"user" yaml:"user"`
	Password  string `mapstructure:"password" yaml:"password"`
	DataId    string `mapstructure:"data_id" yaml:"data_id"`
	Group     string `mapstructure:"group" yaml:"group"`
	Namespace string `mapstructure:"namespace" yaml:"namespace"`
}
