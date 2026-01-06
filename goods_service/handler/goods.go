package handler

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/olivere/elastic/v7"
	"github.com/opentracing/opentracing-go"
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
	"strconv"
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
	var response proto.GoodsInfoResponse
	if goods.ShipFree != nil {
		response.ShipFree = goods.ShipFree
	}
	if goods.IsNew != nil {
		response.IsNew = goods.IsNew
	}
	if goods.IsHot != nil {
		response.IsHot = goods.IsHot
	}
	if goods.OnSale != nil {
		response.OnSale = goods.OnSale
	}
	response.Id = goods.ID
	response.Name = goods.Name
	response.CategoryId = goods.CategoryID
	response.GoodsSn = goods.GoodsSn
	response.ClickNum = goods.ClickNum
	response.SoldNum = goods.SoldNum
	response.FavNum = goods.FavNum
	response.MarketPrice = goods.MarketPrice
	response.ShopPrice = goods.ShopPrice
	response.GoodsBrief = goods.GoodsBrief
	response.GoodsFrontImage = firstImage
	response.DescImages = descImages
	response.Images = otherImages
	response.Brand = &proto.BrandInfoResponse{
		Id:   goods.Brands.ID,
		Name: goods.Brands.Name,
		Logo: goods.Brands.Logo,
	}
	response.Category = &proto.CategoryBriefInfoResponse{
		Id:   goods.Category.ID,
		Name: goods.Category.Name,
	}
	return response
}

func (g GoodSever) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	prepareSpan := opentracing.GlobalTracer().StartSpan("prepare_option", opentracing.ChildOf(parentSpan.Context()))
	response := proto.GoodsListResponse{}
	pageInfo := common.PageInfo{
		Page:  request.Pages,
		Limit: request.PagePerNums,
		Key:   request.KeyWords,
	}
	// 针对 商品的 特殊查询  es 进行查询

	query := elastic.NewBoolQuery()
	if request.IsHot { //是否热卖
		// 这样不加权重  只有模糊匹配 才加权重
		query = query.Filter(elastic.NewTermQuery("is_hot", request.IsHot))
	}
	if request.IsNew { //是否新品
		query = query.Filter(elastic.NewTermQuery("is_new", request.IsNew))
	}
	if request.PriceMax > 0 { //价格区间
		query = query.Filter(elastic.NewRangeQuery("market_price").Lte(request.PriceMax))
	}
	if request.PriceMin > 0 { //价格区间
		query = query.Filter(elastic.NewRangeQuery("market_price").Gte(request.PriceMin))

	}
	if request.BrandID != 0 { // 是否规定品牌
		query = query.Filter(elastic.NewTermQuery("brand_id", request.BrandID))
	}
	if request.KeyWords != "" {
		query = query.Must(elastic.NewMultiMatchQuery(request.KeyWords, "name", "desc"))
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
		type result struct {
			Id int32 `json:"id"`
		}
		var results []result
		search := fmt.Sprintf("category_id in (%s)", subQuery)
		err = global.DB.Model(models.CategoryModel{}).Raw(search).Scan(&results).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Errorf(codes.NotFound, "不存在")
		}
		categoryIds := make([]interface{}, 0)
		for _, v := range results {
			categoryIds = append(categoryIds, strconv.Itoa(int(v.Id)))
		}

		query = query.Filter(elastic.NewTermsQuery("category_id", categoryIds...))

	}
	prepareSpan.Finish()
	// 链路记录
	esSpan := opentracing.GlobalTracer().StartSpan("good_es_search", opentracing.ChildOf(parentSpan.Context()))
	resp, err := global.EsClient.Search().Index(models.EsGoods{}.Index()).Query(query).From(int(pageInfo.GetOffset())).Size(int(pageInfo.GetLimit())).Do(context.Background())
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "查询错误")
	}
	response.Total = int32(resp.Hits.TotalHits.Value)
	esSpan.Finish() //es 完成
	ids := make([]int, 0)
	for _, hit := range resp.Hits.Hits {
		id, err := strconv.Atoi(hit.Id)
		if err != nil {
			zap.S().Error(err)
			return nil, status.Errorf(codes.Internal, "错误")
		}
		ids = append(ids, id)
	}
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("good_mysql_search", opentracing.ChildOf(parentSpan.Context()))
	var list []models.GoodModel
	global.DB.Preload("Category").Preload("Brands").Where("id in (?)", ids).Find(&list)

	var InfoList []*proto.GoodsInfoResponse
	mysqlSpan.Finish() // mysql 查询完成
	for _, item := range list {
		info := GoodInfoFunction(item)
		InfoList = append(InfoList, &info)
	}
	response.Data = InfoList
	return &response, nil

}

func (g GoodSever) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_search", opentracing.ChildOf(parentSpan.Context()))
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
	mysqlSpan.Finish()
	return &proto.GoodsListResponse{
		Total: int32(len(GoodsModels)),
		Data:  response,
	}, nil

}

func (g GoodSever) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_search", opentracing.ChildOf(parentSpan.Context()))
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
	var ImagesModels []*models.GoodsImageModel

	err = tx.Create(&models.GoodsImageModel{
		GoodsID:   model.ID,
		ImageURL:  info.GoodsFrontImage,
		Sort:      0,
		IsMain:    true,
		ImageType: enum.MainImageType, //（1=主图，2=详情图，3=其他）
	}).Error
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "创建失败")
	}
	ImagesModels = append(ImagesModels, &models.GoodsImageModel{
		ImageURL: info.GoodsFrontImage,
	})

	for i, image := range info.DescImages {
		err = tx.Debug().Create(&models.GoodsImageModel{
			GoodsID:   model.ID,
			ImageURL:  image,
			Sort:      i + 1,
			IsMain:    true,
			ImageType: enum.DetailImageType, //（1=主图，2=详情图，3=其他）
		}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "创建失败")
		}
		ImagesModels = append(ImagesModels, &models.GoodsImageModel{
			ImageURL: image,
		})
	}

	for i, image := range info.Images {
		err = tx.Create(&models.GoodsImageModel{
			GoodsID:   model.ID,
			ImageURL:  image,
			Sort:      i + 1,
			IsMain:    true,
			ImageType: enum.OtherImageType, //（1=主图，2=详情图，3=其他）
		}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "创建失败")
		}
		ImagesModels = append(ImagesModels, &models.GoodsImageModel{
			ImageURL: image,
		})
	}
	mysqlSpan.Finish()
	model.Category = &category
	model.Brands = &brand
	model.Images = ImagesModels

	goodInfo := GoodInfoFunction(model)
	err = tx.Commit().Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.Internal, "错误")
	}

	return &goodInfo, nil

}

func (g GoodSever) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_search", opentracing.ChildOf(parentSpan.Context()))
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
	mysqlSpan.Finish()
	return new(empty.Empty), tx.Commit().Error
}

func (g GoodSever) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*empty.Empty, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_search", opentracing.ChildOf(parentSpan.Context()))
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
				Sort:      i + 1,
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
				Sort:      i + 1,
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
	mysqlSpan.Finish()
	return new(empty.Empty), tx.Commit().Error

}

func (g GoodSever) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	parentSpan := opentracing.SpanFromContext(ctx)
	// 链路记录
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_search", opentracing.ChildOf(parentSpan.Context()))
	var GoodsModel models.GoodModel
	err := global.DB.Where("id = ?", request.Id).Preload("Category").Preload("Brands").Preload("Images").Take(&GoodsModel).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	mysqlSpan.Finish()
	goodInfo := GoodInfoFunction(GoodsModel)
	return &goodInfo, nil

}
