// Package natsclient nats.go
package natsclient

import (
	"github.com/nats-io/nats.go"
	"log"
)

var nc *nats.Conn

// InitNats Initialize NATS connection
func InitNats(natsUrl ...string) {
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
}

// Publish a message to a specific subject
func Publish(subject, message string) error {
	return nc.Publish(subject, []byte(message))
}

// Subscribe to a subject and handle messages
func Subscribe(subject string, handler nats.MsgHandler) (*nats.Subscription, error) {
	return nc.Subscribe(subject, handler)
}
