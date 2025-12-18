package api

type Api struct {
	OrderApi OrderApi
	CartApi  CartApi
}

var App = Api{}
