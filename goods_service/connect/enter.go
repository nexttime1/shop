package connect

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"goods_service/global"
	"goods_service/proto"
	"goods_service/utils/otgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InventoryConnectService() (proto.InventoryClient, *grpc.ClientConn, error) {
	//你只要导入这个包  就可以 执行 	resolver.Register(&builder{})  注册进去  就由内部管理  根据 tag 去找到对应服务  轮询的实现负载均衡

	// 这个是 consul 的 ip 和 port
	connectAddr := "consul://" +
		global.Config.ConsulInfo.GetAddr() +
		"/" +
		"stock_service" +
		"?wait=14s"

	zap.S().Infof("try connecting to %s ...", connectAddr)
	conn, err := grpc.NewClient(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())), // 设置拦截器  只要调用 grpc 就拦截
	)
	if err != nil {
		zap.S().Errorf("创建 grpc 客户端连接失败：%v", err)

		return nil, nil, err
	}
	client := proto.NewInventoryClient(conn)
	zap.S().Infof("InventoryClient 连接成功")
	return client, conn, err

}
