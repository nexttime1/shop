package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"order_service/global"
	"order_service/models"
	"order_service/proto"
	"order_service/service"
	"order_service/utils/struct_to_map"
)

// CartItemList 查找某个用户的购物车
func (o OrderSever) CartItemList(ctx context.Context, info *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var cartModels []models.ShoppingCartModel
	var response proto.CartItemListResponse
	count := global.DB.Where("user = ?", info.Id).Find(&cartModels).RowsAffected
	response.Total = int32(count)
	var modelsInfo []*proto.ShopCartInfoResponse
	for _, model := range cartModels {
		modelsInfo = append(modelsInfo, &proto.ShopCartInfoResponse{
			Id:      model.ID,
			UserId:  model.User,
			GoodsId: model.Goods,
			Nums:    model.Nums,
			Checked: model.Checked,
		})

	}
	response.Data = modelsInfo
	return &response, nil
}

func (o OrderSever) CreateCartItem(ctx context.Context, request *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	// 添加商品到购物车   如果存在 那就更新Num + 1
	var model models.ShoppingCartModel
	count := global.DB.Where("user = ?", request.UserId).First(&model).RowsAffected
	if count > 0 {
		// 更新一下
		err := global.DB.Model(&model).Update("nums", model.Nums+1).Error
		if err != nil {
			zap.S().Errorf(err.Error())
			return nil, status.Error(codes.Internal, "加入错误")
		}

	} else {
		model = models.ShoppingCartModel{
			User:    request.UserId,
			Goods:   request.GoodsId,
			Nums:    request.Nums,
			Checked: request.Checked,
		}
		err := global.DB.Create(&model).Error
		if err != nil {
			zap.S().Errorf(err.Error())
			return nil, status.Error(codes.Internal, "创建错误")
		}

	}
	return &proto.ShopCartInfoResponse{
		Id:      model.ID,
		UserId:  model.User,
		GoodsId: model.Goods,
		Nums:    model.Nums,
		Checked: model.Checked,
	}, nil

}

func (o OrderSever) UpdateCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	// 更新check 或者 num
	var model models.ShoppingCartModel
	err := global.DB.Where("id = ?", request.Id).Take(&model).Error
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, status.Error(codes.NotFound, "未找到")
	}
	structMap := service.CartUpdateMap{
		Nums:    model.Nums,
		Checked: model.Checked,
	}
	toMap := struct_to_map.StructToMap(structMap)
	err = global.DB.Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, status.Error(codes.Internal, "更新失败")
	}
	return &emptypb.Empty{}, nil

}

func (o OrderSever) DeleteCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	//删除购物车的某个商品
	var model models.ShoppingCartModel
	err := global.DB.Where("id = ?", request.Id).Take(&model).Error
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, status.Error(codes.NotFound, "未找到")
	}
	err = global.DB.Delete(&model).Error
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, status.Error(codes.Internal, "删除失败")
	}
	return &emptypb.Empty{}, nil
}
