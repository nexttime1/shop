package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"goods_web/common/res"
	"goods_web/connect"
	"goods_web/proto"
	"goods_web/service/pms_srv"
)

type PmsApi struct{}

func (PmsApi) ProductAttrListView(c *gin.Context) {
	var req pms_srv.ProductAttrListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	resp, err := client.ListProductAttr(context.WithValue(context.Background(), "ginContext", c), &proto.ProductAttrListRequest{
		Page:  req.Page,
		Limit: req.Limit,
		Key:   req.Key,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, resp.List, resp.Count)
}

func (PmsApi) ProductAttrCreateView(c *gin.Context) {
	var req pms_srv.ProductAttrCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = client.CreateProductAttr(context.WithValue(context.Background(), "ginContext", c), &proto.ProductAttrItem{
		AttrName:  req.AttrName,
		AttrValue: req.AttrValue,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "创建成功")
}

func (PmsApi) SkuStockListView(c *gin.Context) {
	var req pms_srv.SkuStockListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	resp, err := client.ListSkuStock(context.WithValue(context.Background(), "ginContext", c), &proto.SkuStockListRequest{
		Page:      req.Page,
		Limit:     req.Limit,
		ProductId: req.ProductID,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithList(c, resp.List, resp.Count)
}

func (PmsApi) SkuStockCreateView(c *gin.Context) {
	var req pms_srv.SkuStockCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.GoodConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = client.CreateSkuStock(context.WithValue(context.Background(), "ginContext", c), &proto.SkuStockItem{
		ProductId: req.ProductID,
		SkuCode:   req.SkuCode,
		Stock:     req.Stock,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "创建成功")
}
