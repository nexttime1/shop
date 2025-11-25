package core

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"shop_service/global"
	"shop_service/handler"
	"shop_service/proto"
)

func InitRPC() error {
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserSever{})
	lis, err := net.Listen("tcp", global.Config.UserRPC.GetAddr())
	if err != nil {
		return err
	}

	// 健康注册
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	err = server.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
