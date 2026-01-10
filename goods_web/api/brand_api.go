package api

import (
	"context"
	"github.com/gin-gonic/gin"
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
	res.OkWithData(c, brandInfo)

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
	res.OkWithList(c, list.Data, list.Total)

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
	res.OkWithData(c, Info)
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
