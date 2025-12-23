package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"option_web/common/res"
	"option_web/connect"
	"option_web/proto"
	"option_web/service/collection_srv"
	"option_web/utils/jwts"
)

type UserCollectionApi struct {
}

func (UserCollectionApi) CollectionListView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	collectionClient, conn, err := connect.CollectionConnectService()
	if err != nil {
		res.FailWithMsg(c, res.FailServiceCode, "服务启动失败")
		return
	}
	defer conn.Close()
	list, err := collectionClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: claims.UserID,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	if list.Total == 0 {
		res.OkWithList(c, list.Data, 0)
		return
	}
	var idList []int32
	for _, model := range list.Data {
		idList = append(idList, model.GoodsId)
	}

	goodClient, goodCtConn, err := connect.GoodConnectService()
	if err != nil {
		res.FailWithMsg(c, res.FailServiceCode, "服务启动失败")
		return
	}
	defer goodCtConn.Close()
	goodsInfo, err := goodClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: idList,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	var response []collection_srv.CollectionListResponse

	for _, Info := range list.Data {
		for _, goodModel := range goodsInfo.Data {
			if Info.GoodsId == goodModel.Id {
				response = append(response, collection_srv.CollectionListResponse{
					GoodId:    goodModel.Id,
					Name:      goodModel.Name,
					ShopPrice: goodModel.ShopPrice,
				})
				break
			}
		}
	}
	res.OkWithList(c, response, goodsInfo.Total)

}

func (UserCollectionApi) CollectionAddView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var cr collection_srv.CollectionAddRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}

	collectionClient, conn, err := connect.CollectionConnectService()
	if err != nil {
		res.FailWithMsg(c, res.FailServiceCode, "服务启动失败")
		return
	}
	defer conn.Close()
	_, err = collectionClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  claims.UserID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}

	res.OkWithMessage(c, "收藏成功")
}

func (UserCollectionApi) CollectionDeleteView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var cr collection_srv.CollectionIdRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}

	collectionClient, conn, err := connect.CollectionConnectService()
	if err != nil {
		res.FailWithMsg(c, res.FailServiceCode, "服务启动失败")
		return
	}
	defer conn.Close()
	_, err = collectionClient.DeleteUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  claims.UserID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "删除成功")

}

func (UserCollectionApi) CollectionDetailView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var cr collection_srv.CollectionIdRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}

	collectionClient, conn, err := connect.CollectionConnectService()
	if err != nil {
		res.FailWithMsg(c, res.FailServiceCode, "服务启动失败")
		return
	}
	defer conn.Close()
	_, err = collectionClient.GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		UserId:  claims.UserID,
		GoodsId: cr.GoodId,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "存在")

}
