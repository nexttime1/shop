package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"goods_api/common/res"
	"goods_api/connect"
	"goods_api/proto"
	"goods_api/service/banner_srv"
	"strconv"
)

type BannerApi struct {
}

func (BannerApi) GetBannerListView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	ctx := context.WithValue(context.Background(), "ginContext", c)
	list, err := client.BannerList(ctx, &empty.Empty{})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, list.Data, list.Total)

}

func (BannerApi) CreateBannerView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr banner_srv.BannerCreateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	ctx := context.WithValue(context.Background(), "ginContext", c)
	bannerInfo, err := client.CreateBanner(ctx, &proto.BannerRequest{
		Image: cr.Image,
		Index: cr.Index,
		Url:   cr.Url,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithData(c, bannerInfo)

}

func (BannerApi) DeleteBannerView(c *gin.Context) {
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	var cr banner_srv.BannerIdRequest
	err = c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)

		return
	}
	ctx := context.WithValue(context.Background(), "ginContext", c)
	_, err = client.DeleteBanner(ctx, &proto.BannerRequest{
		Id: cr.Id,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "删除成功")

}

func (BannerApi) UpdateBannerView(c *gin.Context) {
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
	var cr banner_srv.BannerUpdateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	ctx := context.WithValue(context.Background(), "ginContext", c)
	_, err = client.UpdateBanner(ctx, &proto.BannerRequest{
		Id:    int32(id),
		Image: cr.Image,
		Index: cr.Index,
		Url:   cr.Url,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "更新成功")

}
