package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"user_web/common/res"
	"user_web/connect"
	"user_web/proto"
	"user_web/service/ums_srv"
)

type UmsApi struct{}

func (UmsApi) RoleListView(c *gin.Context) {
	var req ums_srv.RoleListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.UserOptionConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	resp, err := client.ListRole(context.WithValue(context.Background(), "ginContext", c), &proto.RoleListRequest{
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

func (UmsApi) RoleDetailView(c *gin.Context) {
	var req ums_srv.RoleIdRequest
	if err := c.ShouldBindUri(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.UserOptionConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	resp, err := client.GetRole(context.WithValue(context.Background(), "ginContext", c), &proto.RoleRequest{Id: req.Id})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithData(c, resp)
}

func (UmsApi) RoleCreateView(c *gin.Context) {
	var req ums_srv.RoleUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.UserOptionConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = client.CreateRole(context.WithValue(context.Background(), "ginContext", c), &proto.RoleItem{
		RoleName: req.RoleName,
		MenuIds:  req.MenuIds,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "创建成功")
}

func (UmsApi) RoleUpdateView(c *gin.Context) {
	var uri ums_srv.RoleIdRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	var req ums_srv.RoleUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.UserOptionConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = client.UpdateRole(context.WithValue(context.Background(), "ginContext", c), &proto.RoleItem{
		Id:       uri.Id,
		RoleName: req.RoleName,
		MenuIds:  req.MenuIds,
	})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "更新成功")
}

func (UmsApi) RoleDeleteView(c *gin.Context) {
	var req ums_srv.RoleIdRequest
	if err := c.ShouldBindUri(&req); err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	client, conn, err := connect.UserOptionConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	_, err = client.DeleteRole(context.WithValue(context.Background(), "ginContext", c), &proto.RoleRequest{Id: req.Id})
	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "删除成功")
}
