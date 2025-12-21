package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"
	"go.uber.org/zap"
	"net/http"
	"order_web/common"
	"order_web/common/enum"
	"order_web/common/res"
	"order_web/connect"
	"order_web/global"
	"order_web/proto"
	"order_web/service/order_srv"
	"order_web/utils/jwts"
	"strconv"
)

type OrderApi struct {
}

func (OrderApi) OrderListView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	userId := claims.UserID
	if claims.Role == enum.AdminRole {
		userId = 0
	}

	orderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	list, err := orderClient.OrderList(context.Background(), &proto.OrderFilterRequest{
		PageNum:  cr.Page,
		PageSize: cr.Limit,
		UserId:   userId,
	})

	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	var response []order_srv.OrderListResponse
	for _, model := range list.Data {
		result := order_srv.OrderListResponse{
			Id:      model.Id,
			UserId:  model.UserId,
			OrderSn: model.OrderSn,
			PayType: model.PayType,
			Status:  model.Status,
			Post:    model.Post,
			Total:   model.Total,
			Address: model.Address,
			Name:    model.Name,
			Mobile:  model.Mobile,
		}
		response = append(response, result)

	}

	res.OkWithList(c, response, list.Total)

}

func (OrderApi) OrderCreateView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var cr order_srv.OrderCreateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {

		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	fmt.Println(cr)
	OrderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	orderModel, err := OrderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  claims.UserID,
		Address: cr.Address,
		Name:    cr.Name,
		Mobile:  cr.Mobile,
		Post:    cr.Post,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}

	client, err := alipay.New(global.Config.Alipay.AppId, global.Config.Alipay.PrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(global.Config.Alipay.AliPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.Config.Alipay.NotifyUrl
	p.ReturnURL = global.Config.Alipay.ReturnUrl
	p.Subject = "下次一定_" + orderModel.OrderSn
	p.OutTradeNo = orderModel.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(orderModel.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	result, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Error(err)
		res.FailWithMsg(c, res.FailServiceCode, "生成支付宝url失败")
		return
	}
	response := order_srv.OrderCreateResponse{
		Id:        orderModel.Id,
		AlipayUrl: result.String(),
	}
	res.OkWithData(c, response)

}

func (OrderApi) DeleteOrderView(c *gin.Context) {

}

func (OrderApi) OrderDetailView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var cr order_srv.OrderIdRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	userId := claims.UserID
	if claims.Role == enum.AdminRole {
		userId = 0
	}
	orderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	result, err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		UserId: userId,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	response := order_srv.OrderDetailResponse{
		Id:      result.OrderInfo.Id,
		UserId:  result.OrderInfo.UserId,
		OrderSn: result.OrderInfo.OrderSn,
		PayType: result.OrderInfo.PayType,
		Status:  result.OrderInfo.Status,
		Post:    result.OrderInfo.Post,
		Total:   result.OrderInfo.Total,
		Address: result.OrderInfo.Address,
		Name:    result.OrderInfo.Name,
		Mobile:  result.OrderInfo.Mobile,
	}
	var goodsInfo []order_srv.GoodInfo
	for _, good := range result.Goods {
		info := order_srv.GoodInfo{
			Id:    good.Id,
			Name:  good.GoodsName,
			Image: good.GoodsImage,
			Price: good.GoodsPrice,
			Nums:  good.Nums,
		}
		goodsInfo = append(goodsInfo, info)
	}
	response.GoodInfo = goodsInfo

	client, err := alipay.New(global.Config.Alipay.AppId, global.Config.Alipay.PrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(global.Config.Alipay.AliPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = global.Config.Alipay.NotifyUrl
	p.ReturnURL = global.Config.Alipay.ReturnUrl
	p.Subject = "下次一定_" + response.OrderSn
	p.OutTradeNo = response.OrderSn
	p.TotalAmount = strconv.FormatFloat(float64(response.Total), 'f', 2, 64)
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	rep, err := client.TradePagePay(p)
	if err != nil {
		zap.S().Error(err)
		res.FailWithMsg(c, res.FailServiceCode, "生成支付宝url失败")
		return
	}
	response.AlipayUrl = rep.String()

	res.OkWithData(c, response)

}

func (OrderApi) AlipayCallBackView(c *gin.Context) {

	client, err := alipay.New(global.Config.Alipay.AppId, global.Config.Alipay.PrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(global.Config.Alipay.AliPublicKey)
	if err != nil {
		panic(err)
	}

	notification, err := client.GetTradeNotification(c.Request)
	if err != nil || notification == nil {
		zap.S().Error(err)
		res.FailWithMsg(c, res.FailServiceCode, "")
		return

	}
	OrderClient, conn, err := connect.OrderConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = OrderClient.UpdateOrderStatus(context.Background(), &proto.OrderStatus{
		OrderSn: notification.OutTradeNo,
		Status:  string(notification.TradeStatus),
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	c.String(http.StatusOK, "success")

}
