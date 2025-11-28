package conf

type Code struct {
	ID     string `mapstructure:"id" yaml:"id"`
	Secret string `mapstructure:"secret" yaml:"secret"`
}
