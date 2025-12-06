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
	"goods_service/service"
	"goods_service/utils/struct_to_map"
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
			subQuery = fmt.Sprintf("select id from category_models where parent_category_id = %d", request.TopCategoryID)
		} else {
			subQuery = fmt.Sprintf("%d", request.TopCategoryID)
		}
		search := fmt.Sprintf("category_id in (%s)", subQuery)
		query = query.Where(search)
	}
	list, count, err := common.ListQuery(models.GoodModel{}, common.Options{
		PageInfo: pageInfo,
		Likes:    []string{"name"},
		Preload:  []string{"Category", "Brands", "Images"},
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
	var GoodsModels []models.GoodModel
	global.DB.Where("id in ?", info.Id).Preload("Category").Preload("Brands").Preload("Images").Find(&GoodsModels)
	if len(GoodsModels) != len(info.Id) {
		return nil, status.Errorf(codes.NotFound, "部分商品不存在")
	}
	var response []*proto.GoodsInfoResponse
	for _, good := range GoodsModels {
		goodInfo := GoodInfoFunction(good)
		response = append(response, &goodInfo)
	}

	return &proto.GoodsListResponse{
		Total: int32(len(GoodsModels)),
		Data:  response,
	}, nil

}

func (g GoodSever) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category models.CategoryModel
	err := global.DB.Where("id = ?", info.CategoryId).Take(&category).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.NotFound, "分类不存在")
	}
	var brand models.Brands
	err = global.DB.Where("id = ?", info.Brand).Take(&brand).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 添加商品
	model := models.GoodModel{
		CategoryID:  info.CategoryId,
		BrandsID:    info.Brand,
		ShipFree:    info.ShipFree,
		Name:        info.Name,
		GoodsSn:     info.GoodsSn,
		ClickNum:    0,
		SoldNum:     0,
		FavNum:      0,
		MarketPrice: info.MarketPrice,
		ShopPrice:   info.ShopPrice,
		GoodsBrief:  info.GoodsBrief,
	}
	err = tx.Create(&model).Error
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建失败")
	}

	// web 已经上传了 七牛云 这里就是 url
	// 添加第三章表  图片
	// 主图
	err = tx.Create(&models.GoodsImageModel{
		GoodsID:   model.ID,
		ImageURL:  info.GoodsFrontImage,
		Sort:      0,
		IsMain:    true,
		ImageType: 1, //（1=主图，2=详情图，3=其他）
	}).Error
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建失败")
	}

	for i, image := range info.DescImages {
		err = tx.Create(&models.GoodsImageModel{
			GoodsID:   model.ID,
			ImageURL:  image,
			Sort:      int32(i + 1),
			IsMain:    true,
			ImageType: 2, //（1=主图，2=详情图，3=其他）
		}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "创建失败")
		}
	}
	for i, image := range info.Images {
		err = tx.Create(&models.GoodsImageModel{
			GoodsID:   model.ID,
			ImageURL:  image,
			Sort:      int32(i + 1),
			IsMain:    true,
			ImageType: 3, //（1=主图，2=详情图，3=其他）
		}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "创建失败")
		}
	}
	goodInfo := GoodInfoFunction(model)

	return &goodInfo, nil

}

func (g GoodSever) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	var model models.GoodModel
	err := global.DB.Where("id = ?", info.Id).Take(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.NotFound, "不存在")
	}
	// 先删除 image 表
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = tx.Where("goods_id = ?", info.Id).Delete(&models.GoodsImageModel{}).Error
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "删除失败")
	}

	err = tx.Delete(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "删除失败")
	}
	return new(empty.Empty), tx.Commit().Error
}

func (g GoodSever) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*empty.Empty, error) {
	fmt.Println("UpdateGoods")
	var model models.GoodModel
	err := global.DB.Where("id = ?", info.Id).Take(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.NotFound, "不存在")
	}
	if info.Brand != 0 {
		var brand models.Brands
		err := global.DB.Where("id = ?", info.Brand).Take(&brand).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Errorf(codes.NotFound, "品牌不存在")
		}

	}
	if info.CategoryId != 0 {
		var category models.CategoryModel
		err := global.DB.Where("id = ?", info.CategoryId).Take(&category).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Errorf(codes.NotFound, "分类不存在")
		}

	}

	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 修改 第三章表
	if info.GoodsFrontImage != "" {
		err = tx.Where("goods_id = ? and is_main = 1", info.Id).Delete(&models.GoodsImageModel{}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "修改错误")
		}
		err = tx.Create(&models.GoodsImageModel{
			GoodsID:   model.ID,
			ImageURL:  info.GoodsFrontImage,
			Sort:      0,
			IsMain:    true,
			ImageType: 1, //（1=主图，2=详情图，3=其他）
		}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "创建失败")
		}
	}
	if info.DescImages != nil {
		err = tx.Where("goods_id = ? and image_type = 2", info.Id).Delete(&models.GoodsImageModel{}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "修改错误")
		}
		for i, image := range info.DescImages {
			err = tx.Create(&models.GoodsImageModel{
				GoodsID:   model.ID,
				ImageURL:  image,
				Sort:      int32(i + 1),
				IsMain:    true,
				ImageType: 2, //（1=主图，2=详情图，3=其他）
			}).Error
			if err != nil {
				zap.S().Error(err)
				tx.Rollback()
				return nil, status.Errorf(codes.Internal, "创建失败")
			}
		}
	}
	if info.Images != nil {
		err = tx.Where("goods_id = ? and image_type = 3", info.Id).Delete(&models.GoodsImageModel{}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "修改错误")
		}

		for i, image := range info.Images {
			err = tx.Create(&models.GoodsImageModel{
				GoodsID:   model.ID,
				ImageURL:  image,
				Sort:      int32(i + 1),
				IsMain:    true,
				ImageType: 3, //（1=主图，2=详情图，3=其他）
			}).Error
			if err != nil {
				zap.S().Error(err)
				tx.Rollback()
				return nil, status.Errorf(codes.Internal, "创建失败")
			}
		}

	}

	// 修改 商品表
	StructMap := service.GoodUpdateServiceMap{
		Name:        info.Name,
		GoodsSn:     info.GoodsSn,
		Stocks:      info.Stocks,
		MarketPrice: info.MarketPrice,
		ShopPrice:   info.ShopPrice,
		GoodsBrief:  info.GoodsBrief,
		ShipFree:    info.ShipFree,
		IsNew:       info.IsNew,
		IsHot:       info.IsHot,
		OnSale:      info.OnSale,
		CategoryId:  info.CategoryId,
		Brand:       info.Brand,
	}

	toMap := struct_to_map.StructToMap(StructMap)
	err = tx.Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "更新错误")
	}
	return new(empty.Empty), tx.Commit().Error

}

func (g GoodSever) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var GoodsModel models.GoodModel
	err := global.DB.Where("id = ?", request.Id).Preload("Category").Preload("Brands").Preload("Images").Take(&GoodsModel).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	goodInfo := GoodInfoFunction(GoodsModel)
	return &goodInfo, nil

}
