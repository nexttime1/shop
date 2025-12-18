package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"order_web/common/res"
	"order_web/connect"
	"order_web/proto"
	"order_web/service/cart_srv"
	"order_web/utils/jwts"
)

type CartApi struct {
}

func (CartApi) CartListView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	OrderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	cartList, err := OrderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: claims.UserID,
	})

	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	idList := make([]int32, 0)
	for _, cart := range cartList.Data {
		idList = append(idList, cart.GoodsId)
	}
	GoodClient, GoodConn, err := connect.GoodConnectService()
	if err != nil {
		return
	}
	defer GoodConn.Close()
	goodsList, err := GoodClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: idList})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	var response []cart_srv.CartListResponse
	for _, item := range cartList.Data {
		for _, goodModel := range goodsList.Data {
			if item.GoodsId == goodModel.Id {
				data := cart_srv.CartListResponse{
					Id:          item.Id,
					GoodID:      item.GoodsId,
					Name:        goodModel.Name,
					GoodsSn:     goodModel.GoodsSn,
					Stocks:      goodModel.Stocks,
					CategoryId:  goodModel.CategoryId,
					MarketPrice: goodModel.MarketPrice,
					GoodPrice:   goodModel.ShopPrice,
					GoodsBrief:  goodModel.GoodsBrief,
					Images:      goodModel.Images,
					DescImages:  goodModel.DescImages,
					ShipFree:    goodModel.ShipFree,
					FrontImage:  goodModel.GoodsFrontImage,
					Chacked:     item.Checked,
				}
				response = append(response, data)
			}
		}
	}
	res.OkWithList(c, response, cartList.Total)

}

func (CartApi) DeleteCartItemView(c *gin.Context) {

}

func (CartApi) AddItemView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var cr cart_srv.CartAddRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	orderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	goodClient, clientConn, err := connect.GoodConnectService()
	if err != nil {
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	defer clientConn.Close()
	goodModel, err := goodClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: cr.GoodID})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}

	// 检查库存
	InventoryClient, inventoryConn, err := connect.InventoryConnectService()
	if err != nil {
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	defer inventoryConn.Close()

	detail, err := InventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: cr.GoodID,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	if cr.Num > detail.Num {
		res.FailWithMsg(c, res.FailArgumentCode, "库存不足")
		return
	}

	checked := true
	req, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:     claims.UserID,
		GoodsId:    cr.GoodID,
		GoodsName:  goodModel.Name,
		GoodsImage: goodModel.GoodsFrontImage,
		GoodsPrice: goodModel.ShopPrice,
		Nums:       cr.Num,
		Checked:    &checked,
	})
	response := cart_srv.CartAddResponse{
		Id: req.Id,
	}

	res.OkWithData(c, response)

}

func (CartApi) UpdatePatchView(c *gin.Context) {

}
