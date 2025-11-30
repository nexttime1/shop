package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"goods_service/proto"
)

type GoodSever struct {
}

func (g GoodSever) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) GetAllCategorysList(ctx context.Context, empty *empty.Empty) (*proto.CategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) BrandList(ctx context.Context, empty *empty.Empty) (*proto.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) BannerList(ctx context.Context, empty *empty.Empty) (*proto.BannerListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}
