package main

import (
	"context"
	"goods_service/models"
	"strconv"

	"goods_service/core"
	"goods_service/flags"
	"goods_service/global"
)

func main() {
	flags.Parse() //解析 yaml文件
	core.InitZap()
	global.Config = core.ReadConf()
	//fmt.Println(global.Config)
	global.DB = core.InitDB()
	global.EsClient = core.InitEs()
	flags.Run()
	mysql2Es()

}

func mysql2Es() {
	var goodModels []models.GoodModel
	global.DB.Find(&goodModels)
	for _, goodModel := range goodModels {
		model := models.EsGoods{
			ID:          goodModel.ID,
			CategoryID:  goodModel.CategoryID,
			BrandsID:    goodModel.BrandsID,
			OnSale:      *goodModel.OnSale,
			ShipFree:    *goodModel.ShipFree,
			IsNew:       *goodModel.IsNew,
			IsHot:       *goodModel.IsHot,
			Name:        goodModel.Name,
			ClickNum:    goodModel.ClickNum,
			SoldNum:     goodModel.SoldNum,
			FavNum:      goodModel.FavNum,
			MarketPrice: goodModel.MarketPrice,
			GoodsBrief:  goodModel.GoodsBrief,
			ShopPrice:   goodModel.ShopPrice,
		}
		_, err := global.EsClient.Index().Index(models.EsGoods{}.Index()).BodyJson(model).Id(strconv.Itoa(int(model.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

}
