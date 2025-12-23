package service

type AddressMap struct {
	Province     string `structs:"province"`
	City         string `structs:"city"`
	District     string `structs:"district"`
	Address      string `structs:"address"`
	SignerName   string `structs:"signer_name"`
	SignerMobile string `structs:"signer_mobile"`
}
