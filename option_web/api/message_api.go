package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"option_web/common/enum"
	"option_web/common/res"
	"option_web/connect"
	"option_web/proto"

	"option_web/service/message_srv"
	"option_web/utils/jwts"
)

type MessageApi struct {
}

func (MessageApi) MessageListView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	request := proto.MessageRequest{
		UserId: claims.UserID,
	}
	if claims.Role == enum.AdminRole {
		request.UserId = 0
	}

	OrderClient, conn, err := connect.MessageConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()
	List, err := OrderClient.MessageList(context.WithValue(context.Background(), "ginContext", c), &request)

	if err != nil {
		zap.S().Error(err)
		res.FailWithServiceMsg(c, err)
		return
	}

	var response []message_srv.MessageResponse

	for _, model := range List.Data {
		response = append(response, message_srv.MessageResponse{
			Id:          model.Id,
			UserId:      model.UserId,
			MessageType: model.MessageType,
			Subject:     model.Subject,
			Message:     model.Message,
			File:        model.File,
		})

	}

	res.OkWithList(c, response, List.Total)

}

func (MessageApi) CreateMessageView(c *gin.Context) {
	_claims, exist := c.Get("claims")
	if !exist {
		return
	}
	claims := _claims.(*jwts.MyClaims)
	var cr message_srv.MessageRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		zap.S().Error(err)
		res.FailWithErr(c, res.FailArgumentCode, err)
		return
	}
	messageClient, conn, err := connect.MessageConnectService(c)
	if err != nil {
		return
	}
	defer conn.Close()

	req, err := messageClient.CreateMessage(context.WithValue(context.Background(), "ginContext", c), &proto.MessageRequest{
		UserId:      claims.UserID,
		MessageType: cr.MessageType,
		Subject:     cr.Subject,
		Message:     cr.Message,
		File:        cr.File,
	})
	if err != nil {
		res.FailWithServiceMsg(c, err)
		return
	}

	res.OkWithData(c, req.Id)

}
