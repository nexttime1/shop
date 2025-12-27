package connect

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"goods_api/common/res"
	"goods_api/global"
	"goods_api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GoodConnectService(c *gin.Context) (proto.GoodsClient, *grpc.ClientConn, error) {
	//你只要导入这个包  就可以 执行 	resolver.Register(&builder{})  注册进去  就由内部管理  根据 tag 去找到对应服务  轮询的实现负载均衡

	// 这个是 consul 的 ip 和 port
	connectAddr := "consul://" +
		global.Config.ConsulInfo.GetAddr() +
		"/goods_service?wait=14s"

	zap.S().Infof("try connecting to %s ...", connectAddr)
	conn, err := grpc.NewClient(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorf("创建 grpc 客户端连接失败：%v", err)
		res.FailWithErr(c, res.FailServiceCode, err)
		return nil, nil, err
	}
	goodClient := proto.NewGoodsClient(conn)
	zap.S().Infof("goodClient 连接成功")
	return goodClient, conn, err

}
func StockConnectService(c *gin.Context) (proto.InventoryClient, *grpc.ClientConn, error) {
	//你只要导入这个包  就可以 执行 	resolver.Register(&builder{})  注册进去  就由内部管理  根据 tag 去找到对应服务  轮询的实现负载均衡

	// 这个是 consul 的 ip 和 port
	connectAddr := "consul://" +
		global.Config.ConsulInfo.GetAddr() +
		"/stock_service?wait=14s"

	zap.S().Infof("try connecting to %s ...", connectAddr)
	conn, err := grpc.NewClient(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Errorf("创建 grpc 客户端连接失败：%v", err)
		res.FailWithErr(c, res.FailServiceCode, err)
		return nil, nil, err
	}
	stockClient := proto.NewInventoryClient(conn)
	zap.S().Infof("stockClient 连接成功")
	return stockClient, conn, err

}

//  演变  请看 User api
