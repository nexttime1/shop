package handler

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_service/global"
	"goods_service/models"
	"goods_service/proto"
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
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) DeleteCategory(ctx context.Context, request *proto.DeleteCategoryRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) UpdateCategory(ctx context.Context, request *proto.CategoryInfoRequest) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}
