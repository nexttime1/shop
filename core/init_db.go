package core

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"shop_service/global"
	"time"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(global.Config.DB.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, //不生成外键约束
	})
	if err != nil {
		logrus.Errorf("数据库连接失败")
		return nil
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("数据库连接成功")
	return db
}
