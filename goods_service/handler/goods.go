package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_service/common"
	"goods_service/global"
	"goods_service/models"
	"goods_service/models/enum"
	"goods_service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GoodInfoFunction(goods models.GoodModel) proto.GoodsInfoResponse {
	// 封面图
	firstImage := ""
	// 详情图 type = 2
	var descImages []string
	// 其他图 type = 3
	var otherImages []string
	for _, imageModel := range goods.Images {
		if imageModel.IsMain {
			firstImage = imageModel.ImageURL
		}
		if imageModel.ImageType == enum.DetailImageType {
			descImages = append(descImages, imageModel.ImageURL)
		}
		if imageModel.ImageType == enum.OtherImageType {
			otherImages = append(otherImages, imageModel.ImageURL)
		}
	}

	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: firstImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      descImages,
		Images:          otherImages,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

func (g GoodSever) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	response := proto.GoodsListResponse{}
	pageInfo := common.PageInfo{
		Page:  request.Pages,
		Limit: request.PagePerNums,
		Key:   request.KeyWords,
	}
	// 针对 商品的 特殊查询
	query := global.DB.Model(models.GoodModel{})
	if request.IsHot { //是否热卖
		query = query.Where("is_hot = true")
	}
	if request.IsNew { //是否新品
		query = query.Where("is_new = true")
	}
	if request.PriceMax > 0 { //价格区间
		query = query.Where("shop_price <= ?", request.PriceMax)
	}
	if request.PriceMin > 0 { //价格区间
		query = query.Where("shop_price >= ?", request.PriceMin)
	}
	if request.BrandID != 0 { // 是否规定品牌
		query = query.Where("brands_id = ?", request.BrandID)
	}
	// 分类的 查询
	var subQuery string
	if request.TopCategoryID != 0 { //说明用户选择了 分类查询
		var model models.CategoryModel
		err := global.DB.Where("id = ?", request.TopCategoryID).Take(&model).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Errorf(codes.NotFound, "没有找到该分类")
		}

		if model.Level == 1 {
			// 一级分类 也就是 找出所有二级分类  	再找到三级分类
			subQuery = fmt.Sprintf("select id from category_models where parent_category_id in (select id from category_models where parent_category_id = %d)", request.TopCategoryID)
		} else if model.Level == 2 {
			subQuery = fmt.Sprintf("select id from category_models where parent_category_id = %d)", request.TopCategoryID)
		} else {
			subQuery = fmt.Sprintf("%d", request.TopCategoryID)
		}
		query = query.Where("category_id in ?", subQuery)
	}
	list, count, err := common.ListQuery(models.GoodModel{}, common.Options{
		PageInfo: pageInfo,
		Likes:    []string{"name"},
		Preload:  []string{"CategoryModel", "Brands", "GoodsImageModel"},
		Where:    query,
		Debug:    true,
	})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "查询错误")
	}
	response.Total = count
	var InfoList []*proto.GoodsInfoResponse
	for _, item := range list {
		info := GoodInfoFunction(item)
		InfoList = append(InfoList, &info)
	}
	response.Data = InfoList
	return &response, nil

}

func (g GoodSever) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}
