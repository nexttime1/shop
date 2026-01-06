package connect

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop_api/common/res"
	"shop_api/global"
	"shop_api/proto"
	"shop_api/utils/otgrpc"
)

func UserConnectService(c *gin.Context) (proto.UserClient, *grpc.ClientConn, error) {
	//你只要导入这个包  就可以 执行 	resolver.Register(&builder{})  注册进去  就由内部管理  根据 tag 去找到对应服务  轮询的实现负载均衡

	// 这个是 consul 的 ip 和 port
	connectAddr := "consul://" +
		global.Config.ConsulInfo.GetAddr() +
		"/user_service?wait=14s"

	zap.S().Infof("try connecting to %s ...", connectAddr)
	conn, err := grpc.NewClient(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), // 设置拦截器  只要调用 grpc 就拦截
	)
	if err != nil {
		zap.S().Errorf("创建 grpc 客户端连接失败：%v", err)
		res.FailWithErr(c, res.FailServiceCode, err)
		return nil, nil, err
	}
	client := proto.NewUserClient(conn)

	return client, conn, err

}

// UserConnectService11 演变 1  已经换新
func UserConnectService11(c *gin.Context) (proto.UserClient, *grpc.ClientConn, error) {
	// 从服务中心去拿
	consulConfig := api.DefaultConfig()
	// 依旧指定consul 在哪
	consulConfig.Address = global.Config.ConsulInfo.GetAddr()
	consulClient, err := api.NewClient(consulConfig)
	// 挑选一个服务
	filterMap, err := consulClient.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.Config.ConsulInfo.Name))
	if err != nil {
		zap.S().Errorf("Consul 未匹配到 %s", err.Error())
		return nil, nil, errors.New("服务错误")
	}
	// 只拿第一个就行
	addr := ""
	port := 0
	for _, service := range filterMap {
		addr = service.Address
		port = service.Port
		break
	}
	if addr == "" {
		return nil, nil, errors.New("服务错误")
	}
	ConnectAddr := fmt.Sprintf("%s:%d", addr, port)

	conn, err := grpc.NewClient(ConnectAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		res.FailWithErr(c, res.FailServiceCode, err)
	}
	client := proto.NewUserClient(conn)

	return client, conn, err
}
