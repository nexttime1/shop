package connect

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop_api/common/res"
	"shop_api/global"
	"shop_api/proto"
)

func UserConnectService(c *gin.Context) (proto.UserClient, *grpc.ClientConn, error) {
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
