package api

type Api struct {
	GoodApi     GoodApi
	CategoryApi CategoryApi
	BannerApi   BannerApi
	BrandApi    BrandApi
	PmsApi      PmsApi
}

var App = Api{}
