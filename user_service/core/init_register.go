package core

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"user_service/handler"
	"user_service/proto"
	"user_service/utils/free_port"
	"user_service/utils/otgrpc"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"user_service/global"
)

type ServerRegister interface {
	Register() error
	Deregister() error
}

type ConsulRegister struct {
	ServiceID string
}

// Go 编译器对指针类型调用值接收者方法时，会自动做 *ptr 解引用
func (c ConsulRegister) Register() error {
	// 动态获得 端口
	port, err := free_port.GetFreePort()
	if err != nil {
		zap.L().Error("端口获得错误 ", zap.Error(err))
		return err
	}
	zap.S().Infof("用户服务获得的端口号为: %d", port)
	server := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())))
	proto.RegisterUserServer(server, &handler.UserSever{})
	proto.RegisterUmsServer(server, &handler.UmsServer{})
	// 监听的端口 一定是动态获取的 要不健康检查 识别不到
	listenAddr := fmt.Sprintf("%s:%d", global.Config.LocalInfo.Addr, port)
	lis, err := net.Listen("tcp", listenAddr)
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
		GRPC:                           fmt.Sprintf("%s:%d", global.Config.LocalInfo.Addr, port), // gRPC 健康检查地址
		GRPCUseTLS:                     false,                                                    // 是否使用 TLS
		Interval:                       "10s",                                                    // 检查间隔
		Timeout:                        "5s",                                                     // 检查超时时间
		DeregisterCriticalServiceAfter: "10s",                                                    // 服务异常后多久注销
	}

	registration := &api.AgentServiceRegistration{
		ID:      c.ServiceID,
		Name:    global.Config.ConsulInfo.Name,
		Tags:    global.Config.ConsulInfo.Tags,
		Address: global.Config.LocalInfo.Addr, // 告诉consul  我这个服务的ip 和 端口
		Port:    port,
		Check:   check,
	}
	//  注册服务到 Consul  发送
	err = consulClient.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Errorf("注册服务到 Consul错误 %s", err)
		return err
	}
	zap.S().Infof("%s 注册成功", global.Config.ConsulInfo.Name)

	// 闭包处理异常
	go func() {
		err := func() error {
			err = server.Serve(lis)
			if err != nil {

				return err
			}
			return nil
		}()
		if err != nil {
			zap.S().Errorf("服务启动失败 %s", err.Error())
			panic(err)
		}
	}()

	return nil
}
func (c ConsulRegister) Deregister() error {
	// 服务注册
	consulConfig := api.DefaultConfig()
	// 改一下默认  这个是 虚拟机Consul 所在的ip
	consulConfig.Address = global.Config.ConsulInfo.GetAddr()

	consulClient, err := api.NewClient(consulConfig)

	err = consulClient.Agent().ServiceDeregister(c.ServiceID)
	if err != nil {
		zap.S().Errorf("服务注销失败 %s", err.Error())
		return err
	}
	return nil

}

func NewConsulRegister() ServerRegister {

	return &ConsulRegister{
		ServiceID: uuid.NewV4().String(),
	}
}
