package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.163.132:9876"}))
	if err != nil {
		panic(err)
	}
	err = p.Start()
	if err != nil {
		panic(err)
	}

	result, err := p.SendSync(context.Background(), primitive.NewMessage("xtm", []byte("hello world")))

	if err != nil {
		panic(err)
	}
	fmt.Println(result.String())
}
