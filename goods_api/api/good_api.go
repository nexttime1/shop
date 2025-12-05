package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goods_api/common/res"
	"goods_api/connect"
	"goods_api/proto"
	"goods_api/service/good_srv"
)

type GoodApi struct {
}

func (GoodApi) GetGoodList(c *gin.Context) {
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
