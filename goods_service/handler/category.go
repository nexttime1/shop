package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_service/global"
	"goods_service/models"
	"goods_service/proto"
	"goods_service/service"
	"goods_service/utils/struct_to_map"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g GoodSever) GetAllCategorysList(ctx context.Context, empty *empty.Empty) (*proto.CategoryListResponse, error) {
	var categoryModels []models.CategoryModel
	global.DB.Debug().Model(&models.CategoryModel{Level: 1}).Preload("SubCategory.SubCategory").Find(&categoryModels)
	bytesData, err := json.Marshal(categoryModels)
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "json marshal error")
	}
	return &proto.CategoryListResponse{
		JsonData: string(bytesData),
	}, nil

}

func (g GoodSever) GetSubCategory(ctx context.Context, request *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var model models.CategoryModel
	err := global.DB.Take(&model, request.Id).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "分类不存在")
	}
	// 如果是一级目录 要显示 到3 级
	preload := "SubCategory"
	if model.Level == 1 {
		preload = "SubCategory.SubCategory"
	}
	var categoryModel models.CategoryModel
	global.DB.Debug().Where("id = ?", model.ID).Preload(preload).Find(&categoryModel)

	response := proto.SubCategoryListResponse{
		Id:             categoryModel.ID,
		Name:           categoryModel.Name,
		ParentCategory: categoryModel.ParentCategoryID,
		IsTab:          categoryModel.IsTab,
		Level:          categoryModel.Level,
	}

	var sub []*proto.SubCategoryListResponse
	for _, c := range categoryModel.SubCategory {
		info := &proto.SubCategoryListResponse{
			Id:             c.ID,
			Name:           c.Name,
			ParentCategory: c.ParentCategoryID,
			IsTab:          c.IsTab,
			Level:          c.Level,
		}
		if c.SubCategory != nil {
			var grandSub []*proto.SubCategoryListResponse
			for _, grandson := range c.SubCategory {
				grandInfo := &proto.SubCategoryListResponse{
					Id:             grandson.ID,
					Name:           grandson.Name,
					ParentCategory: grandson.ParentCategoryID,
					IsTab:          grandson.IsTab,
					Level:          grandson.Level,
				}
				grandSub = append(grandSub, grandInfo)
			}
			info.SubCategories = grandSub
		}
		sub = append(sub, info)
	}
	response.SubCategories = sub

	return &response, nil
}

func (g GoodSever) CreateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	// 看看存不存在
	var model models.CategoryModel
	count := global.DB.Where("name = ?", request.Name).Find(&model).RowsAffected
	if count > 0 {
		return nil, status.Error(codes.AlreadyExists, "分类已经存在")
	}
	if request.Level != 1 {
		// 查一下父分类存不存在
		err := global.DB.Model(&models.CategoryModel{}).Where("id = ?", request.ParentCategoryID).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Error(codes.NotFound, "父分类不存在")
		}
	}
	categoryModel := models.CategoryModel{
		Name:             request.Name,
		ParentCategoryID: request.ParentCategoryID,
		Level:            request.Level,
		IsTab:            request.IsTab,
	}
	err := global.DB.Create(&categoryModel).Error

	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "创建失败")
	}
	return &proto.CategoryInfoResponse{
		Id: categoryModel.ID,
	}, nil
}

func (g GoodSever) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*empty.Empty, error) {
	var model models.CategoryModel
	err := global.DB.Where("id = ?", request.Id).Take(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "分类不存在")
	}
	// 删除中间表
	var brand_category_models []models.BrandCategoryModel
	global.DB.Where("category_id = ?", model.ID).Find(&brand_category_models)
	err = global.DB.Delete(&brand_category_models).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "中间表删除失败")
	}

	// 删除自己
	err = global.DB.Delete(&model).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "删除失败")
	}

	return &empty.Empty{}, nil
}

func (g GoodSever) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*empty.Empty, error) {
	var model models.CategoryModel
	err := global.DB.Take(&model, request.Id).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "分类不存在")
	}
	categoryMap := service.CategoryUpdateServiceMap{
		Name:             request.Name,
		ParentCategoryID: request.ParentCategoryID,
		Level:            request.Level,
		IsTab:            &request.IsTab,
	}
	toMap := struct_to_map.StructToMap(categoryMap)
	err = global.DB.Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "更新失败")
	}
	return &empty.Empty{}, nil

}
