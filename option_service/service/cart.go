package service

type CartUpdateMap struct {
	Nums    int32 `json:"nums" structs:"nums"`
	Checked *bool `json:"checked" structs:"checked"`
}
