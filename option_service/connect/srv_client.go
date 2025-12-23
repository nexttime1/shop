package connect

import (
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"option_service/global"
	"option_service/proto"
)

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
	)
	if err != nil {
		zap.S().Errorf("创建 grpc 客户端连接失败：%v", err)

		return nil, nil, err
	}
	client := proto.NewGoodsClient(conn)
	zap.S().Infof("GoodClient 连接成功")
	return client, conn, err

}
