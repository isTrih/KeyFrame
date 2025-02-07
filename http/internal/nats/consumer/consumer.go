// Package consumer message_consumer.go
package consumer

import (
	"github.com/nats-io/nats.go"
	"log"
	"zerobackend/internal/nats/natsclient"
	"zerobackend/internal/svc"
)

func StartMessageConsumer(svc *svc.ServiceContext) {
	// 订阅 NATS 消息
	_, err := natsclient.Subscribe("message.created", func(msg *nats.Msg) {
		log.Printf("Received a message: %s", string(msg.Data))

		err := svc.BizRedis.Set("nats", string(msg.Data))
		if err != nil {
			return
		}
		// 在这里处理消息，例如插入到数据库、推送到前端等
	})

	if err != nil {
		log.Fatalf("Error subscribing to NATS: %v", err)
	}
	log.Println("Message consumer started")
}
