package api

type Api struct {
	AddressApi    AddressApi
	MessageApi    MessageApi
	CollectionApi CollectionApi
}

var App = Api{}
