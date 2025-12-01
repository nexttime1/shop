package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"goods_service/proto"
)

func (g GoodSever) GoodsList(ctx context.Context, request *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) BatchGetGoods(ctx context.Context, info *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) CreateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) DeleteGoods(ctx context.Context, info *proto.DeleteGoodsInfo) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) UpdateGoods(ctx context.Context, info *proto.CreateGoodsInfo) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (g GoodSever) GetGoodsDetail(ctx context.Context, request *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}
