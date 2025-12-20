package conf

type Alipay struct {
	AppId        string `mapstructure:"app_id" yaml:"app_id"`
	PrivateKey   string `mapstructure:"private_key" yaml:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" yaml:"ali_public_key"`
	NotifyUrl    string `mapstructure:"notify_url" yaml:"notify_url"`
	ReturnUrl    string `mapstructure:"return_url" yaml:"return_url"`
}
