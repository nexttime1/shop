package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"option_service/global"
	"option_service/models"
	"option_service/proto"
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
	//TODO implement me
	panic("implement me")
}

func (o OptionServer) DeleteAddress(ctx context.Context, request *proto.AddressRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (o OptionServer) UpdateAddress(ctx context.Context, request *proto.AddressRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
