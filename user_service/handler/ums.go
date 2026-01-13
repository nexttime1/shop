package handler

import (
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"user_service/common"
	"user_service/global"
	"user_service/models"
	"user_service/proto"
)

type UmsServer struct {
}

func (UmsServer) GetRole(ctx context.Context, req *proto.RoleRequest) (*proto.RoleItem, error) {
	span := opentracing.SpanFromContext(ctx)
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(span.Context()))
	defer mysqlSpan.Finish()

	var role models.Role
	if err := global.DB.Where("id = ?", req.Id).Take(&role).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "role not found")
	}
	menuIDs, _ := parseMenuIDs(role.MenuIDs)
	return &proto.RoleItem{
		Id:       role.ID,
		RoleName: role.RoleName,
		MenuIds:  menuIDs,
	}, nil
}

func (UmsServer) ListRole(ctx context.Context, req *proto.RoleListRequest) (*proto.RoleListResponse, error) {
	span := opentracing.SpanFromContext(ctx)
	mysqlSpan := opentracing.GlobalTracer().StartSpan("mysql_option", opentracing.ChildOf(span.Context()))
	defer mysqlSpan.Finish()

	pageInfo := common.PageInfo{
		Page:  uint32(req.Page),
		Limit: uint32(req.Limit),
		Key:   req.Key,
	}
	options := common.Options{
		PageInfo: pageInfo,
	}
	if req.Key != "" {
		options.Likes = []string{"role_name"}
	}
	list, count, err := common.ListQuery(models.Role{}, options)
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "list role error")
	}
	resp := &proto.RoleListResponse{Count: int32(count)}
	for _, r := range list {
		menuIDs, _ := parseMenuIDs(r.MenuIDs)
		resp.List = append(resp.List, &proto.RoleItem{
			Id:       r.ID,
			RoleName: r.RoleName,
			MenuIds:  menuIDs,
		})
	}
	return resp, nil
}

func (UmsServer) CreateRole(ctx context.Context, req *proto.RoleItem) (*proto.Response, error) {
	menuJSON, err := json.Marshal(req.MenuIds)
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.InvalidArgument, "menu_ids 格式錯誤")
	}
	role := models.Role{
		RoleName: req.RoleName,
		MenuIDs:  menuJSON,
	}
	if err = global.DB.Create(&role).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "create role error")
	}
	return &proto.Response{Code: int32(codes.OK), Msg: "创建成功"}, nil
}

func (UmsServer) UpdateRole(ctx context.Context, req *proto.RoleItem) (*proto.Response, error) {
	menuJSON, err := json.Marshal(req.MenuIds)
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.InvalidArgument, "menu_ids 格式錯誤")
	}
	updates := map[string]interface{}{
		"role_name": req.RoleName,
		"menu_ids":  menuJSON,
	}
	if err = global.DB.Model(&models.Role{}).Where("id = ?", req.Id).Updates(updates).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "update role error")
	}
	return &proto.Response{Code: int32(codes.OK), Msg: "更新成功"}, nil
}

func (UmsServer) DeleteRole(ctx context.Context, req *proto.RoleRequest) (*proto.Response, error) {
	if err := global.DB.Where("id = ?", req.Id).Delete(&models.Role{}).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "delete role error")
	}
	return &proto.Response{Code: int32(codes.OK), Msg: "删除成功"}, nil
}

func parseMenuIDs(data []byte) ([]int32, error) {
	var ids []int32
	err := json.Unmarshal(data, &ids)
	return ids, err
}
