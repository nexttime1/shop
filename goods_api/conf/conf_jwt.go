package conf

type Jwt struct {
	Expire int    `mapstructure:"expire" yaml:"expire"`
	Secret string `mapstructure:"secret" yaml:"secret"`
	Issuer string `mapstructure:"issuer" yaml:"issuer"`
}
