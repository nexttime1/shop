package handler

import (
	"context"
	"go.uber.org/zap"
	"goods_service/common"
	"goods_service/connect"
	"goods_service/global"
	"goods_service/models"
	"goods_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (g GoodSever) GetCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var cb models.BrandCategoryModel
	err := global.DB.Preload("Brands").Preload("Category").Where("id = ?", req.Id).Take(&cb).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "不存在")
	}

	return &proto.CategoryBrandResponse{
		Id: cb.ID,
		Brand: &proto.BrandInfoResponse{
			Id:   cb.Brands.ID,
			Name: cb.Brands.Name,
			Logo: cb.Brands.Logo,
		},
		Category: &proto.CategoryInfoResponse{
			Id:               cb.Category.ID,
			Name:             cb.Category.Name,
			Level:            cb.Category.Level,
			IsTab:            cb.Category.IsTab,
			ParentCategoryID: cb.Category.ParentCategoryID,
		},
	}, nil
}

func (g GoodSever) ListProductAttr(ctx context.Context, req *proto.ProductAttrListRequest) (*proto.ProductAttrListResponse, error) {
	pageInfo := common.PageInfo{
		Page:  req.Page,
		Limit: req.Limit,
	}
	options := common.Options{
		PageInfo: pageInfo,
	}
	if req.Key != "" {
		options.Where = global.DB.Where("attr_name LIKE ?", "%"+req.Key+"%")
	}
	list, count, err := common.ListQuery(models.ProductAttr{}, options)
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "查询失败")
	}

	resp := &proto.ProductAttrListResponse{Count: count}
	for _, item := range list {
		resp.List = append(resp.List, &proto.ProductAttrItem{
			Id:        item.ID,
			AttrName:  item.AttrName,
			AttrValue: item.AttrValue,
		})
	}
	return resp, nil
}

func (g GoodSever) CreateProductAttr(ctx context.Context, req *proto.ProductAttrItem) (*emptypb.Empty, error) {
	attr := models.ProductAttr{
		AttrName:  req.AttrName,
		AttrValue: req.AttrValue,
	}
	if err := global.DB.Create(&attr).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "创建失败")
	}
	return &emptypb.Empty{}, nil
}

func (g GoodSever) ListSkuStock(ctx context.Context, req *proto.SkuStockListRequest) (*proto.SkuStockListResponse, error) {

	if req.ProductId < 0 {
		return nil, status.Error(codes.Internal, "参数错误")
	}
	stockClient, conn, err := connect.InventoryConnectService()
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "链接错误")
	}
	defer conn.Close()
	detail, err := stockClient.InvDetail(ctx, &proto.GoodsInvInfo{
		GoodsId: req.ProductId,
	})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "查询失败")
	}

	resp := &proto.SkuStockListResponse{Count: 1}
	list := &proto.SkuStockItem{
		ProductId: req.ProductId,
		Stock:     detail.Num,
	}
	resp.List = append(resp.List, list)
	return resp, nil
}

func (g GoodSever) CreateSkuStock(ctx context.Context, req *proto.SkuStockItem) (*emptypb.Empty, error) {
	if req.ProductId == 0 {
		return nil, status.Error(codes.InvalidArgument, "商品ID必填")
	}
	var good models.GoodModel
	if err := global.DB.Where("id = ?", req.ProductId).Take(&good).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "商品不存在")
	}

	return &emptypb.Empty{}, nil
}
