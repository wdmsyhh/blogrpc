package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"log"
	"time"
)

func main() {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(primitive.NamesrvAddr{"127.0.0.1:9876"}),
		consumer.WithRetry(2),
		consumer.WithGroupName("GID_a"),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.Subscribe("testmq", consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range ext {
			fmt.Printf("获取到：%+v\n", string(ext[i].Body))
		}
		return consumer.ConsumeSuccess, nil
	})
	err = c.Start()
	defer c.Shutdown()
	time.Sleep(5 * time.Second)
}
