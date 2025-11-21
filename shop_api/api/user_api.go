package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop_api/common/res"
	"shop_api/global"
	"shop_api/proto"
)

type UserApi struct {
}

func (UserApi) UserListView(c *gin.Context) {

	conn, err := grpc.NewClient(global.Config.UserRPC.GetAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		res.FailWithErr(c, res.FailServiceCode, err)
	}
	client := proto.NewUserClient(conn)
	userListResponse, err := client.GetUserList(context.Background(), &proto.PageInfo{
		Page:  1,
		Limit: 5,
	})
	if err != nil {
		res.FailWithMessage(c, err)
		return
	}
	res.OkWithList(c, userListResponse.Data, userListResponse.Total)

}
