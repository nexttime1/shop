package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"goods_service/core"
	"goods_service/flags"
	"goods_service/global"
	"goods_service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn
var client proto.GoodsClient

func Init() {
	flags.Parse() //解析 yaml文件
	global.Config = core.ReadConf()
	global.DB = core.InitDB()
	var err error
	conn, err = grpc.NewClient("192.168.163.1:60120", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewGoodsClient(conn)
}

func GoodList() {
	fmt.Println("GoodList")
	var e *empty.Empty
	list, err := client.BrandList(context.Background(), e)
	if err != nil {
		panic(err)
	}
	for _, brandModel := range list.Data {
		fmt.Println(brandModel)
	}
}

func GetAllCategorys() {
	var e *empty.Empty
	list, err := client.GetAllCategorysList(context.Background(), e)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(list.JsonData)
}

//
//func UserCreate() {
//	for i := 0; i < 3; i++ {
//		request := &proto.CreateUserReq{
//			NickName: fmt.Sprintf("小小test%d", i),
//			Password: "admin",
//			Mobile:   fmt.Sprintf("1756436903%d", i),
//		}
//		user, err := client.CreateUser(context.Background(), request)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(user.Id)
//	}
//
//}
//
//func UserUpdate() {
//	msg, err := client.UpdateUser(context.Background(), &proto.UpdateUserReq{
//		Id:       1,
//		Password: "admin111",
//		NickName: "赵云01",
//		BirthDay: uint64(time.Now().Unix()),
//		Gender:   "male",
//		Role:     1,
//	})
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(msg)
//}

func main() {
	Init()
	defer conn.Close()
	//GoodList()
	GetAllCategorys()
}
