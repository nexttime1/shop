package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"option_web/common/res"
	"option_web/connect"
	"option_web/proto"
	"option_web/service/address_srv"

	"option_web/utils/jwts"
)

type AddressApi struct {
}

func (AddressApi) AddressListView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	addressClient, conn, err := connect.AddressConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	list, err := addressClient.GetAddressList(context.Background(), &proto.AddressRequest{
		UserId: claims.UserID,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	var response []address_srv.AddressListResponse
	for _, model := range list.Data {
		response = append(response, address_srv.AddressListResponse{
			Id:           model.Id,
			UserId:       model.UserId,
			Province:     model.Province,
			City:         model.City,
			District:     model.District,
			Address:      model.Address,
			SignerName:   model.SignerName,
			SignerMobile: model.SignerMobile,
		})
	}
	res.OkWithList(c, response, list.Total)

}

func (AddressApi) AddressCreateView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var cr address_srv.AddressCreateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	AddressClient, conn, err := connect.AddressConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	address, err := AddressClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       claims.UserID,
		Province:     cr.Province,
		City:         cr.City,
		District:     cr.District,
		Address:      cr.Address,
		SignerName:   cr.SignerName,
		SignerMobile: cr.SignerMobile,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	response := address_srv.AddressCreateResponse{
		Id: address.Id,
	}

	res.OkWithData(c, response)

}

func (AddressApi) DeleteAddressView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)

	var cr address_srv.AddressIdRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}

	AddressClient, conn, err := connect.AddressConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = AddressClient.DeleteAddress(context.Background(), &proto.AddressRequest{
		Id:     cr.Id,
		UserId: claims.UserID,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "删除成功")

}

func (AddressApi) UpdateAddressView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var idRequest address_srv.AddressIdRequest
	err := c.ShouldBindUri(&idRequest)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}

	var cr address_srv.AddressUpdateRequest
	err = c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	addressClient, conn, err := connect.AddressConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	_, err = addressClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           idRequest.Id,
		UserId:       claims.UserID,
		Province:     cr.Province,
		City:         cr.City,
		District:     cr.District,
		Address:      cr.Address,
		SignerName:   cr.SignerName,
		SignerMobile: cr.SignerMobile,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}
	res.OkWithMessage(c, "更新成功")

}
