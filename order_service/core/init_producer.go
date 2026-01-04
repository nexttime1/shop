package core

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"order_service/utils/mq"
)

func InitMQ() rocketmq.TransactionProducer {
	var orderListener mq.TransactionProducer

	p, err := rocketmq.NewTransactionProducer(
		&orderListener,
		producer.WithGroupName("order_producer"),
		producer.WithNameServer([]string{"192.168.163.132:9876"}),
	)
	if err != nil {
		zap.S().Error(err)
		return nil
	}

	if err := p.Start(); err != nil {
		zap.S().Error(err)
		return nil
	}

	return p
}
