package flags

import (
	"go.uber.org/zap"
	"goods_service/global"

	"goods_service/models"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.BannerModel{},        // 轮播图
		&models.CategoryModel{},      // 分类表
		&models.Brands{},             //品牌表
		&models.BrandCategoryModel{}, //分类和品牌的 第三个表
		&models.GoodModel{},          //商品表
		&models.GoodsImageModel{},    //商品与 自己图片的关系表
	)
	if err != nil {
		zap.S().Errorf("\n数据库迁移失败  %s", err)
		return
	}
	zap.S().Info("\n数据库迁移成功")

}
