package handler

import (
	"context"
	"github.com/opentracing/opentracing-go"
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
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	var cartModels []models.ShoppingCartModel
	var response proto.CartItemListResponse
	count := global.DB.Where("user = ?", info.Id).Find(&cartModels).RowsAffected
	response.Total = int32(count)
	mysqlSpan.Finish()
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
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	// 添加商品到购物车   如果存在 那就更新Num + 1
	var model models.ShoppingCartModel
	count := global.DB.Where("user = ? and goods = ?", request.UserId, request.GoodsId).First(&model).RowsAffected
	if count > 0 {
		// 更新一下
		err := global.DB.Debug().Model(&model).Update("nums", model.Nums+1).Error
		if err != nil {
			zap.S().Errorf(err.Error())
			return nil, status.Error(codes.Internal, "加入错误")
		}

	} else {
		model = models.ShoppingCartModel{
			User:    request.UserId,
			Goods:   request.GoodsId,
			Nums:    1,
			Checked: request.Checked,
		}
		err := global.DB.Create(&model).Error
		if err != nil {
			zap.S().Errorf(err.Error())
			return nil, status.Error(codes.Internal, "创建错误")
		}
		mysqlSpan.Finish()

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

	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	// 更新check 或者 num
	var model models.ShoppingCartModel
	err := global.DB.Where("user = ? and goods = ?", request.UserId, request.GoodsId).Take(&model).Error
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, status.Error(codes.NotFound, "未找到")
	}
	structMap := service.CartUpdateMap{
		Nums:    request.Nums,
		Checked: request.Checked,
	}
	toMap := struct_to_map.StructToMap(structMap)
	err = global.DB.Debug().Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, status.Error(codes.Internal, "更新失败")
	}
	mysqlSpan.Finish()
	return &emptypb.Empty{}, nil

}

func (o OrderSever) DeleteCartItem(ctx context.Context, request *proto.CartItemRequest) (*emptypb.Empty, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(parentSpan.Context()))
	//删除购物车的某个商品
	var model models.ShoppingCartModel
	err := global.DB.Where("user = ? and goods = ?", request.UserId, request.GoodsId).Take(&model).Error
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, status.Error(codes.NotFound, "未找到")
	}
	err = global.DB.Delete(&model).Error
	if err != nil {
		zap.S().Errorf(err.Error())
		return nil, status.Error(codes.Internal, "删除失败")
	}
	mysqlSpan.Finish()
	return &emptypb.Empty{}, nil
}
