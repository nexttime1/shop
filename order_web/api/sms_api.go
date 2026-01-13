package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"order_web/common/res"
	"order_web/connect"
	"order_web/proto"
	"order_web/service/sms_srv"
	"order_web/utils/jwts"
	"strconv"
)

type SmsApi struct{}

// Coupon
func (SmsApi) CouponListView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var req sms_srv.CouponListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := client.ListCoupon(context.WithValue(context.Background(), "ginContext", c), &proto.CouponListRequest{
		Page:   req.Page,
		Limit:  req.Limit,
		UserID: claims.UserID,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, resp.List, resp.Count)
}

func (SmsApi) CouponDetailView(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		res.FailWithMsg(c, res.FailArgumentCode, "id 参数错误")
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	resp, err := client.GetCoupon(context.WithValue(context.Background(), "ginContext", c), &proto.CouponRequest{Id: int32(id)})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithData(c, resp)
}

func (SmsApi) CouponCreateView(c *gin.Context) {
	var req sms_srv.CouponCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = client.CreateCoupon(context.WithValue(context.Background(), "ginContext", c), &proto.CouponItem{
		CouponCode: req.CouponCode,
		Title:      req.Title,
		Discount:   req.Discount,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "创建成功")
}

// Flash
func (SmsApi) FlashListView(c *gin.Context) {
	var req sms_srv.CouponListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	resp, err := client.ListFlash(context.WithValue(context.Background(), "ginContext", c), &proto.CouponListRequest{
		Page:  req.Page,
		Limit: req.Limit,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, resp.List, resp.Count)
}

func (SmsApi) FlashDetailView(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		res.FailWithMsg(c, res.FailArgumentCode, "id 参数错误")
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	resp, err := client.GetFlash(context.WithValue(context.Background(), "ginContext", c), &proto.CouponRequest{Id: int32(id)})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithData(c, resp)
}

func (SmsApi) FlashCreateView(c *gin.Context) {
	var req sms_srv.FlashCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = client.CreateFlash(context.WithValue(context.Background(), "ginContext", c), &proto.FlashItem{
		Name:    req.Name,
		FlashId: req.FlashID,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "创建成功")
}

// Ad
func (SmsApi) AdListView(c *gin.Context) {
	var req sms_srv.CouponListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	resp, err := client.ListAd(context.WithValue(context.Background(), "ginContext", c), &proto.CouponListRequest{
		Page:  req.Page,
		Limit: req.Limit,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, resp.List, resp.Count)
}

func (SmsApi) AdDetailView(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		res.FailWithMsg(c, res.FailArgumentCode, "id 参数错误")
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	resp, err := client.GetAd(context.WithValue(context.Background(), "ginContext", c), &proto.CouponRequest{Id: int32(id)})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithData(c, resp)
}

func (SmsApi) AdCreateView(c *gin.Context) {
	var req sms_srv.AdCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.OrderSmsConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = client.CreateAd(context.WithValue(context.Background(), "ginContext", c), &proto.AdItem{
		Image: req.Image,
		Url:   req.Url,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "创建成功")
}
