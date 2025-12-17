package main

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"order_service/conf"
	"order_service/core"
	"order_service/flags"
	"order_service/proto"
)

func OrderConnectService() (proto.OrderClient, *grpc.ClientConn, error) {
	//你只要导入这个包  就可以 执行 	resolver.Register(&builder{})  注册进去  就由内部管理  根据 tag 去找到对应服务  轮询的实现负载均衡

	// 这个是 consul 的 ip 和 port
	connectAddr := "consul://" +
		Config.ConsulInfo.GetAddr() +
		"/" +
		"order_service" +
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
	client := proto.NewOrderClient(conn)
	zap.S().Infof("OrderClient 连接成功")
	return client, conn, err

}

var Client proto.OrderClient
var conn *grpc.ClientConn
var err error
var Config *conf.Config

func AddCart() *proto.ShopCartInfoResponse {
	check := true
	response, err2 := Client.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  1,
		GoodsId: 8,
		Checked: &check,
	})
	if err2 != nil {
		zap.S().Error(err2)
		return nil
	}
	return response
}

func CartList() *proto.CartItemListResponse {

	response, err2 := Client.CartItemList(context.Background(), &proto.UserInfo{
		Id: 1,
	})
	if err2 != nil {
		zap.S().Error(err2)
		return nil
	}
	return response
}

func UpdateCheck() error {
	check := false
	_, err2 := Client.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      2,
		Checked: &check,
	})
	if err2 != nil {
		zap.S().Error(err2)
		return nil
	}
	return nil
}

func CreateOrder() error {

	result, err2 := Client.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "北京市",
		Name:    "xtm",
		Mobile:  "18888888888",
		Post:    "尽快发货",
	})
	if err2 != nil {
		zap.S().Error(err2)
		return nil
	}
	fmt.Println(result)
	return nil
}

func main() {
	flags.Parse() //解析 yaml文件
	core.InitZap()
	Config = core.ReadConf()
	Client, conn, err = OrderConnectService()
	if err != nil {
		zap.S().Errorf("错误 %#v", err)
		return
	}
	defer conn.Close()

	//cart := AddCart()
	//fmt.Println(cart.Id)
	//result := CartList()
	//fmt.Println(result)
	//UpdateCheck()
	CreateOrder()

}
