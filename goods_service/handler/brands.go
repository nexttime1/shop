package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_service/common"
	"goods_service/global"
	"goods_service/models"
	"goods_service/proto"
	"goods_service/service"
	"goods_service/utils/struct_to_map"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g GoodSever) BrandList(ctx context.Context, pageInfo *proto.PageInfo) (*proto.BrandListResponse, error) {

	list, count, err := common.ListQuery(models.Brands{}, common.Options{
		PageInfo: common.PageInfo{
			Page:  pageInfo.Page,
			Limit: pageInfo.Limit,
		},
	})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "搜索查询失败")
	}

	// 总返回
	var Response proto.BrandListResponse
	Response.Total = count

	var brandInfo []*proto.BrandInfoResponse
	for _, brandModel := range list {
		brandInfo = append(brandInfo, &proto.BrandInfoResponse{
			Id:   brandModel.ID,
			Name: brandModel.Name,
			Logo: brandModel.Logo,
		})
	}
	Response.Data = brandInfo
	return &Response, nil
}

func (g GoodSever) CreateBrand(ctx context.Context, request *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	// 搜索一下 有没有
	var model models.Brands
	count := global.DB.Where("name = ?", request.Name).Find(&model).RowsAffected
	if count == 1 {
		return nil, status.Error(codes.AlreadyExists, "品牌已经存在")
	}
	model.Name = request.Name
	model.Logo = request.Logo

	err := global.DB.Create(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "数据创建失败")
	}
	response := &proto.BrandInfoResponse{
		Id:   model.ID,
		Name: model.Name,
		Logo: model.Logo,
	}
	return response, nil

}

func (g GoodSever) DeleteBrand(ctx context.Context, request *proto.BrandRequest) (*empty.Empty, error) {

	var model models.Brands
	err := global.DB.Take(&model, request.Id).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 品牌删除  那 商品也要删除
	// 批量获取品牌下的商品ID
	var goodsIDs []int32
	if err = tx.Model(&models.GoodModel{}).
		Where("brands_id = ?", request.Id).
		Pluck("id", &goodsIDs).Error; err != nil {
		tx.Rollback()
		zap.S().Errorf("批量查询品牌[id:%d]下商品ID失败: %v", request.Id, err)
		return nil, status.Error(codes.Internal, "查询商品关联数据失败")
	}

	// 批量删除商品对应的图片
	if len(goodsIDs) > 0 {
		if err = tx.Where("goods_id in ?", goodsIDs).
			Delete(&models.GoodsImageModel{}).Error; err != nil {
			tx.Rollback()
			zap.S().Errorf("批量删除品牌[id:%d]下商品图片失败: %v", request.Id, err)
			return nil, status.Error(codes.Internal, "删除商品图片失败")
		}
		zap.S().Infof("成功删除品牌[id:%d]下%d个商品的图片", request.Id, len(goodsIDs))
	}

	// 批量删除品牌下的所有商品
	if err = tx.Where("brands_id = ?", request.Id).
		Delete(&models.GoodModel{}).Error; err != nil {
		tx.Rollback()
		zap.S().Errorf("批量删除品牌[id:%d]下商品失败: %v", request.Id, err)
		return nil, status.Error(codes.Internal, "删除商品失败")
	}

	err = global.DB.Delete(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "删除失败")
	}
	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback() // 提交失败也需回滚
		zap.S().Errorf("提交删除品牌[id:%d]事务失败: %v", request.Id, err)
		return nil, status.Error(codes.Internal, "事务提交失败")
	}
	return &empty.Empty{}, nil
}

func (g GoodSever) UpdateBrand(ctx context.Context, request *proto.BrandRequest) (*empty.Empty, error) {
	var model models.Brands
	err := global.DB.Take(&model, request.Id).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}
	updateMap := service.BrandUpdateServiceMap{
		Name: request.Name,
		Logo: request.Logo,
	}
	toMap := struct_to_map.StructToMap(updateMap)
	err = global.DB.Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "修改失败")
	}

	return &empty.Empty{}, nil

}
