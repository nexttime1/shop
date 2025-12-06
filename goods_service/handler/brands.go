package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_service/global"
	"goods_service/models"
	"goods_service/proto"
	"goods_service/service"
	"goods_service/utils/struct_to_map"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g GoodSever) BrandList(ctx context.Context, empty *empty.Empty) (*proto.BrandListResponse, error) {
	//TODO 看看需要分页么
	var brandModels []models.Brands
	var count int64
	global.DB.Find(&brandModels).Count(&count)
	// 总返回
	var Response proto.BrandListResponse
	Response.Total = int32(count)

	var brandInfo []*proto.BrandInfoResponse
	for _, brandModel := range brandModels {
		brandInfo = append(brandInfo, &proto.BrandInfoResponse{
			Id:   brandModel.ID,
			Name: brandModel.Name,
			Logo: brandModel.Logo,
		})
	}
	Response.Data = brandInfo
	return &Response, nil
}

func (g GoodSever) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//TODO implement me
	// 搜索一下 有没有
	var model models.Brands
	count := global.DB.Where("name = ?", request.Name).Find(&model).RowsAffected
	if count == 1 {
		return nil, status.Error(codes.AlreadyExists, "品牌已经存在")
	}
	model.Name = request.Name
	model.Logo = request.Logo

	err := global.DB.Create(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "数据创建失败")
	}
	response := &proto.BrandInfoResponse{
		Id:   model.ID,
		Name: model.Name,
		Logo: model.Logo,
	}
	return response, nil

}

func (g GoodSever) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*empty.Empty, error) {
	var model models.Brands
	err := global.DB.Take(&model, request.Id).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	err = global.DB.Delete(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "删除失败")
	}
	return &empty.Empty{}, nil
}

func (g GoodSever) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*empty.Empty, error) {
	var model models.Brands
	err := global.DB.Take(&model, request.Id).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	updateMap := service.BrandUpdateServiceMap{
		Name: request.Name,
		Logo: request.Logo,
	}
	toMap := struct_to_map.StructToMap(updateMap)
	err = global.DB.Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "修改失败")
	}

	return &empty.Empty{}, nil

}
