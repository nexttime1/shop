package handler

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"order_service/common"
	"order_service/global"
	"order_service/models"
	"order_service/proto"
)

type SmsServer struct {
}

// Coupon
func (SmsServer) GetCoupon(ctx context.Context, req *proto.CouponRequest) (*proto.CouponItem, error) {
	var c models.Coupon
	if err := global.DB.Where("id = ?", req.Id).Take(&c).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "coupon not found")
	}
	return &proto.CouponItem{
		Id:         c.ID,
		CouponCode: c.CouponCode,
		Title:      c.Title,
		Discount:   int32(c.Discount),
	}, nil
}

func (SmsServer) ListCoupon(ctx context.Context, req *proto.CouponListRequest) (*proto.CouponListResponse, error) {
	pageInfo := common.PageInfo{
		Page:  req.Page,
		Limit: req.Limit,
	}
	list, count, err := common.ListQuery(models.Coupon{}, common.Options{PageInfo: pageInfo})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "list coupon error")
	}
	resp := &proto.CouponListResponse{Count: int32(count)}
	for _, v := range list {
		resp.List = append(resp.List, &proto.CouponItem{
			Id:         v.ID,
			CouponCode: v.CouponCode,
			Title:      v.Title,
			Discount:   int32(v.Discount),
		})
	}
	return resp, nil
}

func (SmsServer) CreateCoupon(ctx context.Context, req *proto.CouponItem) (*emptypb.Empty, error) {
	c := models.Coupon{
		CouponCode: req.CouponCode,
		Title:      req.Title,
		Discount:   int(req.Discount),
	}
	if err := global.DB.Create(&c).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "create coupon error")
	}
	return &emptypb.Empty{}, nil
}

// Flash
func (SmsServer) GetFlash(ctx context.Context, req *proto.CouponRequest) (*proto.FlashItem, error) {
	var f models.Flash
	if err := global.DB.Where("id = ?", req.Id).Take(&f).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "flash not found")
	}
	return &proto.FlashItem{
		Id:      f.ID,
		Name:    f.Name,
		FlashId: int32(f.FlashID),
	}, nil
}

func (SmsServer) ListFlash(ctx context.Context, req *proto.CouponListRequest) (*proto.CouponListResponse, error) {
	pageInfo := common.PageInfo{
		Page:  req.Page,
		Limit: req.Limit,
	}
	list, count, err := common.ListQuery(models.Flash{}, common.Options{PageInfo: pageInfo})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "list flash error")
	}
	resp := &proto.CouponListResponse{Count: int32(count)}
	for _, v := range list {
		resp.List = append(resp.List, &proto.CouponItem{
			Id:         v.ID,
			CouponCode: "",
			Title:      v.Name,
			Discount:   0,
		})
	}
	return resp, nil
}

func (SmsServer) CreateFlash(ctx context.Context, req *proto.FlashItem) (*emptypb.Empty, error) {
	f := models.Flash{
		Name:    req.Name,
		FlashID: int(req.FlashId),
	}
	if err := global.DB.Create(&f).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "create flash error")
	}
	return &emptypb.Empty{}, nil
}

// Ad
func (SmsServer) GetAd(ctx context.Context, req *proto.CouponRequest) (*proto.AdItem, error) {
	var ad models.Ad
	if err := global.DB.Where("id = ?", req.Id).Take(&ad).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.NotFound, "ad not found")
	}
	return &proto.AdItem{
		Id:    ad.ID,
		Image: ad.Image,
		Url:   ad.Url,
	}, nil
}

func (SmsServer) ListAd(ctx context.Context, req *proto.CouponListRequest) (*proto.CouponListResponse, error) {
	pageInfo := common.PageInfo{
		Page:  req.Page,
		Limit: req.Limit,
	}
	list, count, err := common.ListQuery(models.Ad{}, common.Options{PageInfo: pageInfo})
	if err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "list ad error")
	}
	resp := &proto.CouponListResponse{Count: int32(count)}
	for _, v := range list {
		resp.List = append(resp.List, &proto.CouponItem{
			Id:         v.ID,
			CouponCode: "",
			Title:      "",
			Discount:   0,
		})
	}
	return resp, nil
}

func (SmsServer) CreateAd(ctx context.Context, req *proto.AdItem) (*emptypb.Empty, error) {
	ad := models.Ad{
		Image: req.Image,
		Url:   req.Url,
	}
	if err := global.DB.Create(&ad).Error; err != nil {
		zap.S().Error(err)
		return nil, status.Error(codes.Internal, "create ad error")
	}
	return &emptypb.Empty{}, nil
}
