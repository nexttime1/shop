package api

type Api struct {
	GoodApi     GoodApi
	CategoryApi CategoryApi
	BannerApi   BannerApi
	BrandApi    BrandApi
}

var App = Api{}
