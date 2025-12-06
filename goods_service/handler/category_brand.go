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

func (g GoodSever) CategoryBrandList(ctx context.Context, request *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	// 定义响应结构体
	categoryBrandListResponse := &proto.CategoryBrandListResponse{}

	// 1. 查询总记录数
	pageInfo := common.PageInfo{
		Limit: request.PagePerNums,
		Page:  request.Pages,
	}
	list, count, err := common.ListQuery(models.BrandCategoryModel{}, common.Options{
		PageInfo: pageInfo,
		Preload:  []string{"Category", "Brands"},
	})
	if err != nil {
		zap.S().Error(err)
		return categoryBrandListResponse, status.Error(codes.Internal, "查询失败")
	}

	categoryBrandListResponse.Total = count

	//  转换数据库模型到Proto响应模型
	var categoryResponses []*proto.CategoryBrandResponse
	for _, categoryBrand := range list {
		// 构造分类信息响应
		categoryInfo := &proto.CategoryInfoResponse{
			Id:               categoryBrand.Category.ID,
			Name:             categoryBrand.Category.Name,
			Level:            categoryBrand.Category.Level,
			IsTab:            categoryBrand.Category.IsTab,
			ParentCategoryID: categoryBrand.Category.ParentCategoryID,
		}

		// 构造品牌信息响应
		brandInfo := &proto.BrandInfoResponse{
			Id:   categoryBrand.Brands.ID,
			Name: categoryBrand.Brands.Name,
			Logo: categoryBrand.Brands.Logo,
		}

		// 构造单个分类品牌响应并添加到切片
		categoryResponses = append(categoryResponses, &proto.CategoryBrandResponse{
			Category: categoryInfo,
			Brand:    brandInfo,
		})
	}

	// 赋值响应数据并返回
	categoryBrandListResponse.Data = categoryResponses
	return categoryBrandListResponse, nil
}

// GetCategoryBrandList 一个分类下的所有品牌
func (g GoodSever) GetCategoryBrandList(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	// 最后的返回
	brandListResponse := proto.BrandListResponse{}

	var category models.CategoryModel
	err := global.DB.Where("id = ?", request.Id).Take(&category).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "分类不存在")
	}

	var categoryBrands []models.BrandCategoryModel

	count := global.DB.Preload("Brands").Where(&models.BrandCategoryModel{CategoryID: category.ID}).Find(&categoryBrands).RowsAffected
	brandListResponse.Total = int32(count)

	var brandInfoResponses []*proto.BrandInfoResponse
	for _, categoryBrand := range categoryBrands {
		brandInfoResponses = append(brandInfoResponses, &proto.BrandInfoResponse{
			Id:   categoryBrand.Brands.ID,
			Name: categoryBrand.Brands.Name,
			Logo: categoryBrand.Brands.Logo,
		})
	}

	brandListResponse.Data = brandInfoResponses

	return &brandListResponse, nil
}

func (g GoodSever) CreateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	var category models.CategoryModel
	err := global.DB.Where("id = ?", request.CategoryId).Take(&category).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "分类不存在")
	}

	var brand models.Brands
	err = global.DB.Where("id = ?", request.BrandId).Take(&brand).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "品牌不存在")
	}

	categoryBrand := models.BrandCategoryModel{
		CategoryID: request.CategoryId,
		BrandsID:   request.BrandId,
	}
	err = global.DB.Create(&categoryBrand).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "创建失败")
	}
	response := &proto.CategoryBrandResponse{
		Id: categoryBrand.ID,
		Brand: &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		},
		Category: &proto.CategoryInfoResponse{
			Id:               category.ID,
			Name:             category.Name,
			Level:            category.Level,
			IsTab:            category.IsTab,
			ParentCategoryID: category.ParentCategoryID,
		},
	}

	return response, nil

}

func (g GoodSever) DeleteCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*empty.Empty, error) {
	var categoryBrand models.BrandCategoryModel
	err := global.DB.Where("id = ?", request.Id).Take(&categoryBrand).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "不存在")
	}
	err = global.DB.Delete(&categoryBrand).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "删除错误")
	}
	return &empty.Empty{}, nil
}

func (g GoodSever) UpdateCategoryBrand(ctx context.Context, request *proto.CategoryBrandRequest) (*empty.Empty, error) {
	var categoryBrand models.BrandCategoryModel
	err := global.DB.Where("id = ?", request.Id).Take(&categoryBrand).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "不存在")

	}
	if request.BrandId != 0 {
		var brand models.Brands
		err := global.DB.Where("id = ?", request.BrandId).Take(&brand).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Error(codes.NotFound, "品牌不存在")

		}
	}
	if request.CategoryId != 0 {
		var category models.CategoryModel
		err := global.DB.Where("id = ?", request.CategoryId).Take(&category).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Error(codes.NotFound, "分类不存在")

		}
	}
	ModelMap := service.CategoryBrandUpdateServiceMap{
		CategoryId: request.CategoryId,
		BrandId:    request.BrandId,
	}
	toMap := struct_to_map.StructToMap(ModelMap)

	err = global.DB.Model(&categoryBrand).Updates(toMap).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "更新失败")
	}
	return &empty.Empty{}, nil

}
