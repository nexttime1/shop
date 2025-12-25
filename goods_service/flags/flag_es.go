package flags

import (
	"go.uber.org/zap"
	"goods_service/models"
)

func FlagES() {
	//err := models.ArticleModel{}.CreateIndex()
	//if err != nil {
	//	logrus.Error(err)
	//}
	err := models.EsGoods{}.CreateIndex()
	if err != nil {
		zap.S().Error(err)
	}
}
