package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goods_api/common"
	"goods_api/common/res"
	"goods_api/connect"
	"goods_api/proto"
	"goods_api/service/brand_srv"
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
	list, err := client.BrandList(context.Background(), &proto.PageInfo{
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
	brandInfo, err := client.CreateBrand(context.Background(), &proto.BrandRequest{
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
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = client.UpdateBrand(context.Background(), &proto.BrandRequest{
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
	_, err = client.DeleteBrand(context.Background(), &proto.BrandRequest{
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

	list, err := client.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
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
	list, err := client.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
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

	Info, err := client.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
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
	_, err = client.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
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
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = client.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
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
