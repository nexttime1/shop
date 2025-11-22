package conf

type Jwt struct {
	Expire int    `mapstructure:"expire"`
	Secret string `mapstructure:"secret"`
	Issuer string `mapstructure:"issuer"`
}
