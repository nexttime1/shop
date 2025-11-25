package core

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
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

	// 健康检查注册 gRPC 服务端 内部注册一个服务   Consul Client 来调用它
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	consulConfig := api.DefaultConfig()
	// 改一下默认  这个是 虚拟机Consul 所在的ip
	consulConfig.Address = global.Config.ConsulInfo.GetAddr()

	consulClient, err := api.NewClient(consulConfig)
	//健康检查配置 用于放到 请求的json中  告诉他如何访问  用rpc 而不是 http 去检查健康
	check := &api.AgentServiceCheck{
		GRPC:                           global.Config.LocalInfo.GetAddr(), // gRPC 健康检查地址
		GRPCUseTLS:                     false,                             // 是否使用 TLS
		Interval:                       "10s",                             // 检查间隔
		Timeout:                        "5s",                              // 检查超时时间
		DeregisterCriticalServiceAfter: "10s",                             // 服务异常后多久注销
	}
	// 服务注册请求体 申请注册表
	registration := &api.AgentServiceRegistration{
		ID:      global.Config.ConsulInfo.Name,
		Name:    global.Config.ConsulInfo.Name,
		Tags:    []string{"xtm", "skw", "love"},
		Address: global.Config.LocalInfo.Addr, // 告诉consul  我这个服务的ip 和 端口
		Port:    global.Config.LocalInfo.Port,
		Check:   check,
	}
	//  注册服务到 Consul  发送
	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorf("注册服务到 Consul错误 %s", err.Error())
		return err
	}
	zap.S().Infof("%s 注册成功", global.Config.ConsulInfo.Name)
	err = server.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
