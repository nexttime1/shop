package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"time"
)

func main() {
	messgaes, err := rocketmq.NewPushConsumer(consumer.WithNameServer([]string{"192.168.163.132:9876"}), consumer.WithGroupName("shop"))
	if err != nil {
		panic(err)
	}
	messgaes.Start()
	defer messgaes.Shutdown()
	messgaes.Subscribe("xtm", consumer.MessageSelector{}, func(ctx context.Context, m ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range m {
			fmt.Printf("消息是：%v", m[i].Body)
		}
		return consumer.ConsumeSuccess, nil
	})
	time.Sleep(time.Hour)

}
