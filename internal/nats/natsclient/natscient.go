// Package natsclient nats.go
package natsclient

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"time"
)

var nc *nats.Conn

// InitNats Initialize NATS connection
func InitNats(natsUrl ...string) {
	// In the `jetstream` package, almost all API calls rely on `context.Context` for timeout/cancellation handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var url string
	if len(natsUrl) > 0 {
		url = natsUrl[0]
	} else {
		url = nats.DefaultURL
	}

	var err error
	nc, err = nats.Connect(url)

	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	log.Println("Connected to NATS server")
	// JetStream实例
	js, _ := jetstream.New(nc)

	// 创建stream
	_, _ = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:              "MESSAGES",
		Subjects:          []string{"KEYFRAME.MSG.*"},
		Storage:           jetstream.FileStorage,
		MaxBytes:          1024 * 1024 * 1024,
		MaxMsgs:           1000000,
		NoAck:             false,
		Discard:           jetstream.DiscardOld,
		Retention:         jetstream.LimitsPolicy,
		MaxAge:            265 * 24 * time.Hour,
		MaxMsgsPerSubject: 1000000,
		MaxMsgSize:        1024 * 1024 * 1024,
		MaxConsumers:      1000000,
	})
}

// Publish a message to a specific subject
func Publish(subject, message string) error {
	return nc.Publish(subject, []byte(message))
}

// Subscribe to a subject and handle messages
func Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return nc.Subscribe(subject, handler)
}
