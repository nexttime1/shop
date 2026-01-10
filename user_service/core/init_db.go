package core

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
	"user_service/global"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(global.Config.DB.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //不生成外键约束
	})
	if err != nil {
		zap.S().Errorf("数据库连接失败")
		return nil
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	zap.S().Infof("数据库连接成功")
	return db
}
