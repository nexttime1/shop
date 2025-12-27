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
	var response []good_srv.GoodsInfoResponse
	for _, model := range list.Data {
		info := good_srv.GoodsInfoResponse{
			ID:              model.Id,
			CategoryID:      model.CategoryId,
			Name:            model.Name,
			GoodsSn:         model.GoodsSn,
			ClickNum:        model.ClickNum,
			SoldNum:         model.SoldNum,
			FavNum:          model.FavNum,
			Stocks:          model.Stocks,
			MarketPrice:     model.MarketPrice,
			ShopPrice:       model.ShopPrice,
			GoodsBrief:      model.GoodsBrief,
			GoodsDesc:       model.GoodsDesc,
			ShipFree:        model.ShipFree,
			Images:          model.Images,
			DescImages:      model.DescImages,
			GoodsFrontImage: model.GoodsFrontImage,
			IsNew:           model.IsNew,
			IsHot:           model.IsHot,
			OnSale:          model.OnSale,
			AddTime:         model.AddTime,
			Category: good_srv.CategoryBriefInfoResponse{
				ID:   model.Category.Id,
				Name: model.Category.Name,
			},
			Brand: good_srv.BrandInfoResponse{
				ID:   model.Brand.Id,
				Name: model.Brand.Name,
				Logo: model.Brand.Logo,
			},
		}
		response = append(response, info)
	}
	res.OkWithList(c, response, list.Total)

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
	stockClient, clientConn, err := connect.StockConnectService(c)
	if err != nil {
		return
	}
	defer clientConn.Close()
	detail, err := stockClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodInfo.Id,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	response := good_srv.GoodsInfoResponse{
		ID:              goodInfo.Id,
		CategoryID:      goodInfo.CategoryId,
		Name:            goodInfo.Name,
		GoodsSn:         goodInfo.GoodsSn,
		ClickNum:        goodInfo.ClickNum,
		SoldNum:         goodInfo.SoldNum,
		FavNum:          goodInfo.FavNum,
		Stocks:          detail.Num,
		MarketPrice:     goodInfo.MarketPrice,
		ShopPrice:       goodInfo.ShopPrice,
		GoodsBrief:      goodInfo.GoodsBrief,
		GoodsDesc:       goodInfo.GoodsDesc,
		ShipFree:        goodInfo.ShipFree,
		Images:          goodInfo.Images,
		DescImages:      goodInfo.DescImages,
		GoodsFrontImage: goodInfo.GoodsFrontImage,
		IsNew:           goodInfo.IsNew,
		IsHot:           goodInfo.IsHot,
		OnSale:          goodInfo.OnSale,
		AddTime:         goodInfo.AddTime,
		Category: good_srv.CategoryBriefInfoResponse{
			ID:   goodInfo.Category.Id,
			Name: goodInfo.Category.Name,
		},
		Brand: good_srv.BrandInfoResponse{
			ID:   goodInfo.Brand.Id,
			Name: goodInfo.Brand.Name,
			Logo: goodInfo.Brand.Logo,
		},
	}

	res.OkWithData(c, response)

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
