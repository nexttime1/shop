package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_service/global"
	"goods_service/models"
	"goods_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g GoodSever) GetAllCategorysList(ctx context.Context, empty *empty.Empty) (*proto.CategoryListResponse, error) {
	var categoryModels []models.CategoryModel
	global.DB.Debug().Model(&models.CategoryModel{Level: 1}).Preload("SubCategory.SubCategory").Find(&categoryModels)
	bytesData, err := json.Marshal(categoryModels)
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "json marshal error")
	}
	return &proto.CategoryListResponse{
		JsonData: string(bytesData),
	}, nil

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
