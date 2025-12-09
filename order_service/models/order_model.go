package models

import "time"

// OrderModel 订单主表
type OrderModel struct {
	Model
	User         int32     `gorm:"type:int;index;comment:用户ID"`
	OrderSn      string    `gorm:"type:varchar(30);index;comment:订单编号（唯一）"`
	PayType      string    `gorm:"type:varchar(20);comment:支付方式（alipay/wechat）"`
	Status       string    `gorm:"type:varchar(20);comment:订单状态（PAYING/TRADE_SUCCESS等）"`
	TradeNo      string    `gorm:"type:varchar(100);comment:第三方支付交易号"`
	OrderMount   float32   `gorm:"comment:订单总金额"`
	PayTime      time.Time `gorm:"comment:支付时间"`
	Address      string    `gorm:"type:varchar(100);comment:收货地址"`
	SignerName   string    `gorm:"type:varchar(20);comment:签收人姓名"`
	SignerMobile string    `gorm:"type:varchar(11);comment:签收人手机号"`
	Post         string    `gorm:"type:varchar(20);comment:物流单号"`
}

// TableName 重写订单主表表名
func (OrderModel) TableName() string {
	return "orderinfo"
}
