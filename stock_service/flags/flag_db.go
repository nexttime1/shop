package flags

import (
	"go.uber.org/zap"
	"stock_service/global"
	"stock_service/models"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.InventoryModel{},
	)
	if err != nil {
		zap.S().Errorf("\n数据库迁移失败  %s", err)
		return
	}
	zap.S().Info("\n数据库迁移成功")

}
