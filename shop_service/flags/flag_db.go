package flags

import (
	"go.uber.org/zap"
	"shop_service/global"
	"shop_service/models"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.UserModel{},
	)
	if err != nil {
		zap.S().Errorf("\n数据库迁移失败  %s", err)
		return
	}
	zap.S().Info("\n数据库迁移成功")

}
