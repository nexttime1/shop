package api

type Api struct {
	OrderApi OrderApi
	CartApi  CartApi
	SmsApi   SmsApi
}

var App = Api{}
