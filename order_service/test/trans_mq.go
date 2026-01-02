package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type TransactionProducer struct {
}

// When send transactional prepare(half) message succeed, this method will be invoked to execute local transaction.
func (TransactionProducer) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

// When no response to prepare(half) message. broker will send check message to check the transaction status, and this
// method will be invoked to get local transaction status.
func (TransactionProducer) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	return primitive.CommitMessageState
}

func main() {
	transactionProducer, err := rocketmq.NewTransactionProducer(
		&TransactionProducer{},
		producer.WithNameServer([]string{"192.168.163.132:9876"}),
	)
	if err != nil {
		panic(err)
	}
	err = transactionProducer.Start()
	if err != nil {
		panic(err)
	}
	transaction, err := transactionProducer.SendMessageInTransaction(context.Background(), primitive.NewMessage("skw", []byte("stranction")))
	if err != nil {
		panic(err)
	}
	fmt.Println(transaction.String())

}
