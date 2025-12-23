package address_srv

type AddressListResponse struct {
	Id           int32  `json:"id"`
	UserId       int32  `json:"userId"`
	Province     string `json:"province"`
	City         string `json:"city"`
	District     string `json:"district"`
	Address      string `json:"address"`
	SignerName   string `json:"signer_name"`
	SignerMobile string `json:"signer_mobile"`
}
type AddressCreateRequest struct {
	Province     string `form:"province" json:"province" binding:"required"`
	City         string `form:"city" json:"city" binding:"required"`
	District     string `form:"district" json:"district" binding:"required"`
	Address      string `form:"address" json:"address" binding:"required"`
	SignerName   string `form:"signer_name" json:"signer_name" binding:"required"`
	SignerMobile string `form:"signer_mobile" json:"signer_mobile" binding:"required"`
}

type AddressCreateResponse struct {
	Id           int32  `json:"id"`
	UserId       int32  `json:"userId"`
	Province     string `json:"province"`
	City         string `json:"city"`
	District     string `json:"district"`
	Address      string `json:"address"`
	SignerName   string `json:"signer_name"`
	SignerMobile string `json:"signer_mobile"`
}
type AddressIdRequest struct {
	Id int32 `uri:"id" binding:"required,min=1"`
}

type AddressUpdateRequest struct {
	Province     string `json:"province"`
	City         string `json:"city"`
	District     string `json:"district"`
	Address      string `json:"address"`
	SignerName   string `json:"signer_name"`
	SignerMobile string `json:"signer_mobile"`
}

type OrderListResponse struct {
	Id      int32   `json:"id"`
	UserId  int32   `json:"user_id"`
	OrderSn string  `json:"order_sn"`
	PayType string  `json:"pay_type"`
	Status  string  `json:"status"`
	Post    string  `json:"post"`
	Total   float32 `json:"total"`
	Address string  `json:"address"`
	Name    string  `json:"name"`
	Mobile  string  `json:"mobile"`
}

type OrderDetailResponse struct {
	Id        int32      `json:"id"`
	UserId    int32      `json:"user_id"`
	OrderSn   string     `json:"order_sn"`
	PayType   string     `json:"pay_type"`
	Status    string     `json:"status"`
	Post      string     `json:"post"`
	Total     float32    `json:"total"`
	Address   string     `json:"address"`
	Name      string     `json:"name"`
	Mobile    string     `json:"mobile"`
	GoodInfo  []GoodInfo `json:"goods"`
	AlipayUrl string     `json:"alipay_url"`
}

type GoodInfo struct {
	Id    int32   `json:"id"`
	Name  string  `json:"name"`
	Image string  `json:"image"`
	Price float32 `json:"price"`
	Nums  int32   `json:"nums"`
}
