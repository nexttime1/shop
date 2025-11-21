package core

import (
	"google.golang.org/grpc"
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
	err = server.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
