package handler

import (
	"context"
	"encoding/json"
	"fmt"
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
		Total:    int32(len(categoryModels)),
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
		var parentCategory models.CategoryModel
		err := global.DB.Model(&models.CategoryModel{}).Where("id = ?", request.ParentCategoryID).Take(&parentCategory).Error
		if err != nil {
			zap.S().Error(err)
			return nil, status.Error(codes.NotFound, "父分类不存在")
		}
	}
	categoryModel := models.CategoryModel{
		Name:             request.Name,
		ParentCategoryID: request.ParentCategoryID,
		Level:            request.Level,
	}
	if request.IsTab != nil {
		categoryModel.IsTab = *request.IsTab
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
	// 开启事务，保证操作原子性
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除中间表
	err = tx.Where("category_id = ?", model.ID).Delete(&models.BrandCategoryModel{}).Error
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		return nil, status.Error(codes.Internal, "中间表删除失败")
	}
	// 删除自己的子分类
	subQuery := ""
	if model.Level == 1 {
		subQuery = fmt.Sprintf("select id from category_models where parent_category_id in (select id from category_models where parent_category_id = %d)", model.ID)
	}
	if model.Level == 2 {
		subQuery = fmt.Sprintf("select id from category_models where parent_category_id = %d", model.ID)
	}
	if model.Level != 3 {
		err = tx.Where(fmt.Sprintf("category_id in (%s)", subQuery)).Delete(&models.BrandCategoryModel{}).Error
		if err != nil {
			zap.S().Error(err)
			tx.Rollback()
			return nil, status.Error(codes.Internal, "删除失败")
		}
	}

	// 删除自己
	err = tx.Delete(&model).Error
	if err != nil {
		zap.S().Error(err)
		tx.Rollback()
		return nil, status.Error(codes.Internal, "删除失败")
	}

	if err = tx.Commit().Error; err != nil {
		zap.S().Errorf("事务提交失败：%v", err)
		tx.Rollback()
		return nil, status.Error(codes.Internal, "事务提交失败")
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
		Name:  request.Name,
		IsTab: request.IsTab,
	}
	toMap := struct_to_map.StructToMap(categoryMap)
	err = global.DB.Model(&model).Updates(toMap).Error
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "更新失败")
	}
	return &empty.Empty{}, nil

}
