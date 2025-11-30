package models

import "goods_service/models/enum"

type GoodsImageModel struct {
	Model
	GoodsID   int32          `gorm:"type:int;not null;comment:商品ID（逻辑外键，关联good_models.id）;index:idx_goods_image_goods"`
	ImageURL  string         `gorm:"type:varchar(255);not null;comment:图片访问URL（七牛云）"`
	Sort      int32          `gorm:"type:int;not null;default:0;comment:排序序号（越小越靠前）"`
	IsMain    bool           `gorm:"default:false;not null;comment:是否主图（一个商品仅一个主图）"`
	ImageType enum.ImageType `gorm:"type:tinyint(1);not null;default:3;comment:图片类型（1=主图，2=详情图，3=其他）"`
}
