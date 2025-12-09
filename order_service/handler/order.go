package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"order_service/common"
	"order_service/core"
	"order_service/global"
	"order_service/models"
	"order_service/proto"
)

type OrderSever struct {
}

func (o OrderSever) CreateOrder(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	//先拿到 选中的 good ID
	check := true
	var goodsId []int32
	var shopModels []models.ShoppingCartModel
	global.DB.Where(models.ShoppingCartModel{
		User:    request.UserId,
		Checked: &check,
	}).Find(&shopModels)
	goodNumMap := make(map[int32]int32)
	for _, shopModel := range shopModels {
		goodsId  = append(goodsId, shopModel.Goods)
		goodNumMap[shopModel.Goods] = shopModel.Nums
	}


	// 调用good 微服务
	goodClient, conn, err := core.GoodConnectService()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建失败"))
	}
	defer conn.Close()
	goods, err := goodClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: goodsId,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "商品查询失败")
	}
	var PriceSum float32
	var orderGoods []*models.OrderGoodsModel
	for _, goodModel := range goods.Data {
		PriceSum += goodModel.ShopPrice * float32(goodNumMap[goodModel.Id])
		orderGoods = append(orderGoods, &models.OrderGoodsModel{
			Goods : goodModel.Id,
			GoodsName : goodModel.Name,
			GoodsPrice : goodModel.ShopPrice,
			GoodImages: goodModel.GoodsFrontImage,
			Nums : goodNumMap[goodModel.Id],
		})
	}


}

func (o OrderSever) OrderList(ctx context.Context, request *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	// 管理员看所有的列表   而用户看自己的  区别是 看web端给我发不发id
	var response *proto.OrderListResponse
	pageInfo := common.PageInfo{
		Page:  request.PageNum,
		Limit: request.PageSize,
	}
	var err error
	var list []models.OrderModel
	var count int
	if request.UserId == 0 {
		// 管理员看所有
		list, count, err = common.ListQuery(models.OrderModel{}, common.Options{
			PageInfo: pageInfo,
		})
		if err != nil {
			zap.S().Error(err)
			return nil, status.Error(codes.Internal, "查询错误")
		}

	} else {
		list, count, err = common.ListQuery(models.OrderModel{User: request.UserId}, common.Options{
			PageInfo: pageInfo,
		})
		if err != nil {
			zap.S().Error(err)
			return nil, status.Error(codes.Internal, "查询错误")
		}
	}
	response.Total = int32(count)
	var modelsInfo []*proto.OrderInfoResponse
	for _, item := range list {
		modelsInfo = append(modelsInfo, &proto.OrderInfoResponse{
			Id:      item.ID,
			UserId:  item.User,
			OrderSn: item.OrderSn,
			PayType: item.PayType,
			Status:  item.Status,
			Post:    item.Post,
			Total:   item.OrderMount,
			Address: item.Address,
			Name:    item.SignerName,
			Mobile:  item.SignerMobile,
		})
	}
	response.Data = modelsInfo
	return response, nil
}

func (o OrderSever) OrderDetail(ctx context.Context, request *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	// 如果传userId  那就查这个用户的  不传就是全部的
	var response *proto.OrderInfoDetailResponse
	var model models.OrderModel
	err := global.DB.Where(models.OrderModel{User: request.UserId, Model: models.Model{ID: request.Id}}).Take(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "订单不存在")
	}
	response.OrderInfo = &proto.OrderInfoResponse{
		Id:      model.ID,
		UserId:  model.User,
		OrderSn: model.OrderSn,
		PayType: model.PayType,
		Status:  model.Status,
		Post:    model.Post,
		Total:   model.OrderMount,
		Address: model.Address,
		Name:    model.SignerName,
		Mobile:  model.SignerMobile,
	}
	// 找一下商品
	var goodModels []models.OrderGoodsModel
	global.DB.Where("order = ?", model.ID).Find(&goodModels)
	var Goods []*proto.OrderItemResponse
	for _, item := range goodModels {
		Goods = append(Goods, &proto.OrderItemResponse{
			Id:         item.ID,
			OrderId:    item.Order,
			GoodsId:    item.Goods,
			GoodsName:  item.GoodsName,
			GoodsPrice: item.GoodsPrice,
			Nums:       item.Nums,
		})
	}
	response.Goods = Goods
	return response, nil

}

func (o OrderSever) UpdateOrderStatus(ctx context.Context, status *proto.OrderStatus) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
