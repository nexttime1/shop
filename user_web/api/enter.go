package api

type Api struct {
	UserApi    UserApi
	CaptchaApi CaptchaApi
	UmsApi     UmsApi
}

var App = Api{}
