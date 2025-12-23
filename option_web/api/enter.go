package api

type Api struct {
	AddressApi    AddressApi
	CollectionApi UserCollectionApi
	MessageApi    MessageApi
}

var App = Api{}
