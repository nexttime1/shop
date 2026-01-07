package models

type OrderStockHistory struct {
	Model
	OrderSn string `json:"order_sn"`
	Status  uint   `json:"status" comment:"0 代表 没扣减库存 1 代表 扣减库存"`
}
