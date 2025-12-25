package core

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"goods_service/global"
	"log"
	"os"
)

func InitEs() *elastic.Client {
	//创建客户端实例   把嗅探功能关闭
	es := global.Config.EsInfo
	logger := log.New(os.Stdout, "xtm_", log.LstdFlags)
	url := fmt.Sprintf("%s:%d", es.Addr, es.Port)
	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("http://%s", url)),
		elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}
	zap.S().Infof("ES 连接成功")
	return client
}
