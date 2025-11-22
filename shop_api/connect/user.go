package connect

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop_api/common/res"
	"shop_api/global"
	"shop_api/proto"
)

func UserConnectService(c *gin.Context) (proto.UserClient, *grpc.ClientConn, error) {

	conn, err := grpc.NewClient(global.Config.UserRPC.GetAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		res.FailWithErr(c, res.FailServiceCode, err)
	}
	client := proto.NewUserClient(conn)

	return client, conn, err
}
