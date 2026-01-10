package api

type Api struct {
	UserApi    UserApi
	CaptchaApi CaptchaApi
}

var App = Api{}
