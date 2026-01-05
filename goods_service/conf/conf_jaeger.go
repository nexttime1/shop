package conf

type JaegerInfo struct {
	Host string `mapstructure:"host" yaml:"host"`
	Port uint64 `mapstructure:"port" yaml:"port"`
}
