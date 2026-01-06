package connect

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"order_web/common/res"
	"order_web/global"
	"order_web/proto"
	"order_web/utils/otgrpc"
)

func OrderConnectService(c *gin.Context) (proto.OrderClient, *grpc.ClientConn, error) {
	//你只要导入这个包  就可以 执行 	resolver.Register(&builder{})  注册进去  就由内部管理  根据 tag 去找到对应服务  轮询的实现负载均衡

	// 这个是 consul 的 ip 和 port
	connectAddr := "consul://" +
		global.Config.ConsulInfo.GetAddr() +
		"/order_service?wait=14s"

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
	client := proto.NewOrderClient(conn)
	zap.S().Infof("Client 连接成功")
	return client, conn, err

}

func GoodConnectService() (proto.GoodsClient, *grpc.ClientConn, error) {
	//你只要导入这个包  就可以 执行 	resolver.Register(&builder{})  注册进去  就由内部管理  根据 tag 去找到对应服务  轮询的实现负载均衡

	// 这个是 consul 的 ip 和 port
	connectAddr := "consul://" +
		global.Config.ConsulInfo.GetAddr() +
		"/" +
		global.Config.GoodSrv.Name +
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
	client := proto.NewGoodsClient(conn)
	zap.S().Infof("GoodClient 连接成功")
	return client, conn, err

}

func InventoryConnectService() (proto.InventoryClient, *grpc.ClientConn, error) {
	//你只要导入这个包  就可以 执行 	resolver.Register(&builder{})  注册进去  就由内部管理  根据 tag 去找到对应服务  轮询的实现负载均衡

	// 这个是 consul 的 ip 和 port
	connectAddr := "consul://" +
		global.Config.ConsulInfo.GetAddr() +
		"/" +
		global.Config.InventorySrv.Name +
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
