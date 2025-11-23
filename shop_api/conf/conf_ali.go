package conf

type Code struct {
	ID     string `mapstructure:"id"`
	Secret string `mapstructure:"secret"`
}
