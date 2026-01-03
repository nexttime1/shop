package service

type OrderTransitionRequest struct {
	Id       int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`          // 订单ID（查询详情时必填，创建时不传）
	UserId   int32  `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`  // 用户ID
	Address  string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"` // 收货地址
	Name     string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`       // 签收人姓名
	Mobile   string `protobuf:"bytes,5,opt,name=mobile,proto3" json:"mobile,omitempty"`   // 签收人手机号
	Post     string `protobuf:"bytes,6,opt,name=post,proto3" json:"post,omitempty"`       // 物流单号（创建时可选，发货后填充）
	OrderSns string
}
