package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"shop_service/core"
	"shop_service/flags"
	"shop_service/global"
	"shop_service/proto"
	"time"
)

var conn *grpc.ClientConn
var client proto.UserClient

func Init() {
	flags.Parse() //解析 yaml文件
	global.Config = core.ReadConf()
	global.DB = core.InitDB()
	core.InitLogrus()
	var err error
	conn, err = grpc.NewClient(global.Config.System.GetAddr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewUserClient(conn)
}

func UserList() {
	pageInfo := &proto.PageInfo{
		Page:  1,
		Limit: 2,
	}
	list, err := client.GetUserList(context.Background(), pageInfo)
	if err != nil {
		panic(err)
	}
	for _, user := range list.Data {
		fmt.Println(user)
	}
}

func UserCreate() {
	for i := 0; i < 3; i++ {
		request := &proto.CreateUserReq{
			NickName: fmt.Sprintf("小小test%d", i),
			Password: "admin",
			Mobile:   fmt.Sprintf("1756436903%d", i),
		}
		user, err := client.CreateUser(context.Background(), request)
		if err != nil {
			panic(err)
		}
		fmt.Println(user.Id)
	}

}

func UserUpdate() {
	msg, err := client.UpdateUser(context.Background(), &proto.UpdateUserReq{
		Id:       1,
		Password: "admin111",
		NickName: "赵云01",
		BirthDay: uint64(time.Now().Unix()),
		Gender:   "male",
		Role:     1,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)
}

func main() {
	Init()
	defer conn.Close()
	//UserList()
	//UserCreate()
	UserUpdate()

}
