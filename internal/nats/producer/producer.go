// Package producer.go
package producer

import (
	"log"
	"zerobackend/internal/nats/natsclient"
)

// SendMessageToQueue Send a message to the NATS queue
func SendMessageToQueue(message string) error {
	// Publish message to the "message.created" topic
	err := natsclient.Publish("message.created", message)
	if err != nil {
		log.Printf("Failed to send message to NATS: %v", err)
		return err
	}
	log.Printf("Message sent to NATS: %s", message)
	return nil
}
