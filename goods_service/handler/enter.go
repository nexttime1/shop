package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"goods_service/proto"
)

type GoodSever struct {
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
