package conf

type Log struct {
	App string `mapstructure:"app" yaml:"app"`
	Dir string `mapstructure:"dir" yaml:"dir"`
}
