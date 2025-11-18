package flags

import (
	"github.com/sirupsen/logrus"
	"shop_service/global"
	"shop_service/models"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.UserModel{},
	)
	if err != nil {
		logrus.Errorf("\n数据库迁移失败  %s", err)
		return
	}
	logrus.Info("\n数据库迁移成功")

}
