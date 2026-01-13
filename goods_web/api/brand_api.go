package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"goods_web/common"
	"goods_web/common/res"
	"goods_web/connect"
	"goods_web/proto"
	"goods_web/service/brand_srv"

	"strconv"
)

type BrandApi struct{}

func (BrandApi) BrandListView(c *gin.Context) {
	var cr common.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	list, err := client.BrandList(ctx, &proto.PageInfo{
		Page:  cr.Page,
		Limit: cr.Limit,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, list.Data, list.Total)

}

func (BrandApi) CreateBrandView(c *gin.Context) {
	var cr brand_srv.BrandCreateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	brandInfo, err := client.CreateBrand(ctx, &proto.BrandRequest{
		Name: cr.Name,
		Logo: cr.Logo,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	RMap := map[string]interface{}{
		"id": brandInfo.Id,
	}
	res.OkWithData(c, RMap)

}

func (BrandApi) UpdateBrandView(c *gin.Context) {
	var cr brand_srv.BrandUpdateRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}

	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	_, err = client.UpdateBrand(ctx, &proto.BrandRequest{
		Id:   int32(id),
		Name: cr.Name,
		Logo: cr.Logo,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "修改成功")

}

func (BrandApi) DeleteBrandView(c *gin.Context) {
	var cr brand_srv.BrandIdRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	_, err = client.DeleteBrand(ctx, &proto.BrandRequest{
		Id: cr.Id,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "删除成功")

}

func (BrandApi) BrandDetailView(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "BrandDetailView")
	defer span.Finish()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		res.FailWithMsg(c, res.FailArgumentCode, "id 参数错误")
		return
	}

	Client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := Client.GetBrand(ctx, &proto.BrandRequest{Id: int32(id)})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}

	res.OkWithData(c, resp)
}

//第三张表

func (BrandApi) CategoryBrandListView(c *gin.Context) {
	var cr common.PageInfo
	if err := c.ShouldBindQuery(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	list, err := client.CategoryBrandList(ctx, &proto.CategoryBrandFilterRequest{
		Pages:       cr.Page,
		PagePerNums: cr.Limit,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	var response []brand_srv.BrandCategoryItem
	for _, model := range list.Data {
		response = append(response, brand_srv.BrandCategoryItem{
			Brand: brand_srv.Brand{
				Id:   model.Brand.Id,
				Name: model.Brand.Name,
				Logo: model.Brand.Logo,
			},

			Category: brand_srv.Category{
				Id:               model.Category.Id,
				Name:             model.Category.Name,
				ParentCategoryID: model.Category.ParentCategoryID,
				Level:            model.Category.Level,
				IsTab:            model.Category.IsTab,
			},
		})
	}

	res.OkWithList(c, response, list.Total)

}

func (BrandApi) CategoryAllBrandView(c *gin.Context) {
	var cr brand_srv.BrandIdRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	list, err := client.GetCategoryBrandList(ctx, &proto.CategoryInfoRequest{
		Id: cr.Id,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, list.Data, list.Total)

}

func (BrandApi) CreateCategoryBrandView(c *gin.Context) {
	var cr brand_srv.CreateCategoryBrandRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	Info, err := client.CreateCategoryBrand(ctx, &proto.CategoryBrandRequest{
		BrandId:    cr.BrandId,
		CategoryId: cr.CategoryId,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	RMap := map[string]interface{}{
		"id": Info.Id,
	}
	res.OkWithData(c, RMap)
}

func (BrandApi) DeleteCategoryBrandView(c *gin.Context) {
	var cr brand_srv.BrandIdRequest
	if err := c.ShouldBindUri(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	_, err = client.DeleteCategoryBrand(ctx, &proto.CategoryBrandRequest{
		Id: cr.Id,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "删除成功")

}

func (BrandApi) UpdateCategoryBrandView(c *gin.Context) {
	var cr brand_srv.UpdateCategoryBrandRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}

	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	_, err = client.UpdateCategoryBrand(ctx, &proto.CategoryBrandRequest{
		Id:         int32(id),
		BrandId:    cr.BrandId,
		CategoryId: cr.CategoryId,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "更新成功")

}
