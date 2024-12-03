package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	Reader *kafka.Reader
}

func (c *Consumer) StartConsumer(ctx context.Context, processMessage func(event BookingEvent) error) error {
	for {
		// Читаем сообщения
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Преобразуемм сообщения в структуру BookingEvent
		var event BookingEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		if err := processMessage(event); err != nil {
			log.Printf("Failed to process message: %v", err)
			continue
		}

		if err := c.Reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Failed to commit message: %v", err)
			continue
		}
		log.Printf("Message committed successfully: %s", string(msg.Value))
	}
}
