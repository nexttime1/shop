package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"option_service/global"
	"option_service/models"
	"option_service/proto"
	"option_service/service"
	"option_service/utils/struct_to_map"
)

func AddressFunction(model models.Address) proto.AddressResponse {
	response := proto.AddressResponse{
		Id:           model.ID,
		UserId:       model.UserId,
		Province:     model.Province,
		City:         model.City,
		District:     model.District,
		Address:      model.Address,
		SignerName:   model.SignerName,
		SignerMobile: model.SignerMobile,
	}
	return response
}

func (o OptionServer) GetAddressList(ctx context.Context, request *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var response proto.AddressListResponse
	var addresses []*proto.AddressResponse
	var addressModels []models.Address
	global.DB.Where("user_id = ?", request.UserId).Find(&addressModels)
	for _, model := range addressModels {
		result := AddressFunction(model)
		addresses = append(addresses, &result)
	}

	response.Total = int32(len(addresses))
	response.Data = addresses
	return &response, nil

}

func (o OptionServer) CreateAddress(ctx context.Context, request *proto.AddressRequest) (*proto.AddressResponse, error) {
	model := models.Address{
		UserId:       request.UserId,
		Province:     request.Province,
		City:         request.City,
		District:     request.District,
		Address:      request.Address,
		SignerName:   request.SignerName,
		SignerMobile: request.SignerMobile,
	}
	err := global.DB.Create(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "创建失败")
	}
	response := proto.AddressResponse{
		Id: model.ID,
	}
	return &response, nil

}

func (o OptionServer) DeleteAddress(ctx context.Context, request *proto.AddressRequest) (*emptypb.Empty, error) {
	var model models.Address
	err := global.DB.Where("id = ? and user_id = ?", request.Id, request.UserId).Take(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "不存在")
	}
	return &emptypb.Empty{}, nil

}

func (o OptionServer) UpdateAddress(ctx context.Context, request *proto.AddressRequest) (*emptypb.Empty, error) {
	var model models.Address
	err := global.DB.Where("id = ? and user_id = ?", request.Id, request.UserId).Take(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "不存在")
	}
	modelMap := service.AddressMap{
		Province:     model.Province,
		City:         model.City,
		District:     model.District,
		Address:      model.Address,
		SignerName:   model.SignerName,
		SignerMobile: model.SignerMobile,
	}
	toMap := struct_to_map.StructToMap(modelMap)
	err = global.DB.Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "更新失败")
	}

	response := emptypb.Empty{}
	return &response, nil

}
