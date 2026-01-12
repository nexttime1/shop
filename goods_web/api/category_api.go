package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_web/common/res"
	"goods_web/connect"
	"goods_web/proto"
	"goods_web/service/category_srv"
	"strconv"
)

type CategoryApi struct {
}

func (CategoryApi) GetAllCategoryView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	list, err := client.GetAllCategorysList(context.WithValue(context.Background(), "ginContext", c), &empty.Empty{})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	var response []interface{}

	err = json.Unmarshal([]byte(list.JsonData), &response)
	if err != nil {
		zap.S().Error(err)
		res.FailWithMsg(c, res.FailServiceCode, "json 解析错误")
		return
	}

	res.OkWithList(c, response, list.Total)

}

func (CategoryApi) GetSubCategoryView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr category_srv.CategoryIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	categoryInfo, err := client.GetSubCategory(context.WithValue(context.Background(), "ginContext", c), &proto.CategoryListRequest{
		Id: cr.Id,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	response := ProtoToWebSubCategory(categoryInfo)
	res.OkWithData(c, response)

}

func ProtoToWebSubCategory(protoCategory *proto.SubCategoryListResponse) *category_srv.SubCategoryResponse {
	if protoCategory == nil {
		return nil
	}
	// 1. 赋值当前层级的分类基础字段
	webCategory := &category_srv.SubCategoryResponse{
		Id:             protoCategory.Id,
		Name:           protoCategory.Name,
		ParentCategory: protoCategory.ParentCategory,
		Level:          protoCategory.Level,
		IsTab:          protoCategory.IsTab,
	}

	// 2. 递归处理【子分类】，完美适配 1→2→3 级嵌套
	if protoCategory.SubCategories != nil && len(protoCategory.SubCategories) > 0 {
		webSubList := make([]*category_srv.SubCategoryResponse, 0, len(protoCategory.SubCategories))
		for _, protoSub := range protoCategory.SubCategories {
			webSub := ProtoToWebSubCategory(protoSub)
			webSubList = append(webSubList, webSub)
		}
		webCategory.SubCategories = webSubList
	}
	return webCategory
}

func (CategoryApi) CreateCategoryView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr category_srv.CategoryCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	category, err := client.CreateCategory(context.WithValue(context.Background(), "ginContext", c), &proto.CategoryInfoRequest{
		Name:             cr.Name,
		ParentCategoryID: cr.ParentCategory,
		Level:            cr.Level,
		IsTab:            cr.IsTab,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	RMap := map[string]interface{}{
		"id": category.Id,
	}
	res.OkWithData(c, RMap)

}

func (CategoryApi) UpdateCategoryView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}

	var cr category_srv.UpdateCategoryRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	_, err = client.UpdateCategory(context.WithValue(context.Background(), "ginContext", c), &proto.CategoryInfoRequest{
		Id:    int32(id),
		Name:  cr.Name,
		IsTab: cr.IsTab,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	res.OkWithMessage(c, "更新成功")

}

func (CategoryApi) DeleteCategoryView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr category_srv.CategoryIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	_, err = client.DeleteCategory(context.WithValue(context.Background(), "ginContext", c), &proto.DeleteCategoryRequest{
		Id: cr.Id,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	res.OkWithMessage(c, "删除成功")

}
