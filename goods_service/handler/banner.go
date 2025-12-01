package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_service/global"
	"goods_service/models"
	"goods_service/proto"
	"goods_service/utils/struct_to_map"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func BannerFunction(model models.BannerModel) proto.BannerResponse {
	return proto.BannerResponse{
		Id:    model.ID,
		Image: model.Image,
		Url:   model.Url,
		Index: model.Index,
	}
}

func (g GoodSever) BannerList(ctx context.Context, empty *empty.Empty) (*proto.BannerListResponse, error) {
	var response proto.BannerListResponse
	var bannerModels []models.BannerModel
	count := global.DB.Find(&bannerModels).RowsAffected
	// 如果没有 那就用本地默认
	if count == 0 {
		response.Total = 1
		var bannerList []*proto.BannerResponse
		bannerList = append(bannerList, &proto.BannerResponse{
			Id:    1,
			Image: "http://127.0.0.1:8080/default/xxx.png",
		})
		response.Data = bannerList
		return &response, nil
	}
	response.Total = int32(count)
	var bannerList []*proto.BannerResponse
	for _, model := range bannerModels {
		bannerInfo := BannerFunction(model)
		bannerList = append(bannerList, &bannerInfo)
	}
	response.Data = bannerList

	return &response, nil
}

func (g GoodSever) CreateBanner(ctx context.Context, request *proto.BannerRequest) (*proto.BannerResponse, error) {

	model := models.BannerModel{
		Image: request.Image,
		Url:   request.Url,
		Index: request.Index,
	}

	err := global.DB.Create(&model).Error
	if err != nil {
		zap.S().Error(err.Error())
		return nil, status.Error(codes.Internal, "创建失败")
	}
	return &proto.BannerResponse{
		Id:    model.ID,
		Image: model.Image,
		Url:   model.Url,
		Index: model.Index,
	}, nil

}

func (g GoodSever) DeleteBanner(ctx context.Context, request *proto.BannerRequest) (*empty.Empty, error) {
	var model models.BannerModel
	err := global.DB.Take(&model, request.Id).Error
	if err != nil {
		zap.S().Error(err.Error())
		return nil, status.Error(codes.NotFound, "图片不存在")
	}
	err = global.DB.Delete(&model).Error
	if err != nil {
		zap.S().Error(err.Error())
		return nil, status.Error(codes.Internal, "删除失败")
	}
	return &empty.Empty{}, nil
}

func (g GoodSever) UpdateBanner(ctx context.Context, request *proto.BannerRequest) (*empty.Empty, error) {
	var model models.BannerModel
	err := global.DB.Take(&model, request.Id).Error
	if err != nil {
		zap.S().Error(err.Error())
		return nil, status.Error(codes.NotFound, "图片不存在")
	}
	updateMap := map[string]interface{}{
		"image": request.Image,
		"url":   request.Url,
		"index": request.Index,
	}
	toMap := struct_to_map.StructToMap(updateMap)
	err = global.DB.Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "修改失败")
	}

	return &empty.Empty{}, nil
}
