package models

import (
	"context"
	"go.uber.org/zap"
	"goods_service/global"
)

type EsGoods struct {
	ID          int32   `json:"id"`
	CategoryID  int32   `json:"category_id"`
	BrandsID    int32   `json:"brand_id"`
	OnSale      bool    `json:"on_sale"`
	ShipFree    bool    `json:"ship_free"`
	IsNew       bool    `json:"is_new"`
	IsHot       bool    `json:"is_hot"`
	Name        string  `json:"name"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavNum      int32   `json:"fav_num"`
	MarketPrice float32 `json:"market_price"`
	GoodsBrief  string  `json:"goods_brief"`
	ShopPrice   float32 `json:"shop_price"`
}

func (EsGoods) GetMapping() string {
	goodsMapping := `
{
    "mappings": {
        "properties": {
            "brands_id": {
                "type": "integer"
            },
            "category_id": {
                "type": "integer"
            },
            "click_num": {
                "type": "integer"
            },
            "fav_num": {
                "type": "integer"
            },
            "id": {
                "type": "integer"
            },
            "is_hot": {
                "type": "boolean"
            },
            "is_new": {
                "type": "boolean"
            },
            "market_price": {
                "type": "float"
            },
            "name": {
                "type": "text",
                "analyzer": "ik_max_word"
            },
            "goods_brief": {
                "type": "text",
                "analyzer": "ik_max_word"
            },
            "on_sale": {
                "type": "boolean"
            },
            "ship_free": {
                "type": "boolean"
            },
            "shop_price": {
                "type": "float"
            },
            "sold_num": {
                "type": "long"
            }
        }
    }
}`
	return goodsMapping
}

func (EsGoods) Index() string {
	return "goods_index"
}

// IndexExists 索引是否存在
func (a EsGoods) IndexExists() bool {
	exists, err := global.EsClient.
		IndexExists(a.Index()).
		Do(context.Background())
	if err != nil {
		zap.S().Error(err)
		return exists
	}
	return exists
}

// CreateIndex 创建索引
func (a EsGoods) CreateIndex() error {
	if a.IndexExists() {
		// 有索引
		a.RemoveIndex()
	}
	// 没有索引
	// 创建索引
	createIndex, err := global.EsClient.
		CreateIndex(a.Index()).
		BodyString(a.GetMapping()).
		Do(context.Background())
	if err != nil {
		zap.S().Error("创建索引失败")
		zap.S().Error(err)
		return err
	}
	if !createIndex.Acknowledged {
		zap.S().Error("创建失败")
		return err
	}
	zap.S().Infof("索引 %s 创建成功", a.Index())
	return nil
}

// RemoveIndex 删除索引
func (a EsGoods) RemoveIndex() error {
	zap.S().Info("索引存在，删除索引")
	// 删除索引
	indexDelete, err := global.EsClient.DeleteIndex(a.Index()).Do(context.Background())
	if err != nil {
		zap.S().Error("删除索引失败")
		zap.S().Error(err.Error())
		return err
	}
	if !indexDelete.Acknowledged {
		zap.S().Error("删除索引失败")
		return err
	}
	zap.S().Info("索引删除成功")
	return nil
}
