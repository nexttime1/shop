package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goods_api/common/res"
	"goods_api/connect"
	"goods_api/proto"
	"goods_api/service/good_srv"
	"strconv"
)

type GoodApi struct {
}

func (GoodApi) GetGoodListView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr good_srv.GoodListRequest
	err = c.ShouldBindQuery(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	list, err := client.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		PriceMin:      cr.PriceMin,
		PriceMax:      cr.PriceMax,
		IsHot:         cr.IsHot,
		IsNew:         cr.IsNew,
		TopCategoryID: cr.TopCategoryID,
		Pages:         cr.Page,
		PagePerNums:   cr.Limit,
		KeyWords:      cr.Key,
		BrandID:       cr.BrandID,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, list.Data, list.Total)

}

func (GoodApi) CreateGoodView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr good_srv.GoodCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	goodInfo, err := client.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            cr.Name,
		GoodsSn:         cr.GoodsSn,
		Stocks:          cr.Stocks,
		MarketPrice:     cr.MarketPrice,
		ShopPrice:       cr.ShopPrice,
		GoodsBrief:      cr.GoodsBrief,
		ShipFree:        cr.ShipFree,
		Images:          cr.Images,
		DescImages:      cr.DescImages,
		GoodsFrontImage: cr.FrontImage,
		CategoryId:      cr.CategoryId,
		Brand:           cr.Brand,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithData(c, goodInfo)
}

func (GoodApi) GoodDetailView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr good_srv.GoodDetailRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	goodInfo, err := client.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: cr.Id,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithData(c, goodInfo)

}

func (GoodApi) GoodUpdateView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr good_srv.GoodUpdateRequest
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}

	err = c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	zap.S().Info(cr)

	_, err = client.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(id),
		Name:            cr.Name,
		GoodsSn:         cr.GoodsSn,
		Stocks:          cr.Stocks,
		MarketPrice:     cr.MarketPrice,
		ShopPrice:       cr.ShopPrice,
		GoodsBrief:      cr.GoodsBrief,
		ShipFree:        cr.ShipFree,
		Images:          cr.Images,
		DescImages:      cr.DescImages,
		GoodsFrontImage: cr.FrontImage,
		CategoryId:      cr.CategoryId,
		Brand:           cr.Brand,
	})

	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "更新成功")

}

func (GoodApi) GoodPatchUpdateView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr good_srv.GoodPatchUpdateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	_, err = client.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		IsNew:  cr.IsNew,
		IsHot:  cr.IsHot,
		OnSale: cr.OnSale,
	})

	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "更新成功")
}

func (GoodApi) GoodDeleteView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr good_srv.GoodDeleteRequest

	err = c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	fmt.Println(cr)

	_, err = client.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{
		Id: cr.Id,
	})

	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "删除成功")

}
