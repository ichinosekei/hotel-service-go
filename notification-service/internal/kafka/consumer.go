package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader *kafka.Reader
}

func (c *Consumer) Start(ctx context.Context, processMessage func(event BookingEvent) error) {
	for {
		if ctx.Err() != nil {
			log.Println("Consumer context closed, exiting loop.")
			return
		}

		// Читаем
		msg, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil { // Если контекст завершён, выходим
				log.Println("Consumer context closed, exiting loop.")
				return
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Преобразуем сообщения в структуру BookingEvent
		var event BookingEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// Обрабатываем сообщение
		if err := processMessage(event); err != nil {
			log.Printf("Failed to process message: %v", err)
			continue
		}

		// Подтверждаем обработку сообщения
		if err := c.Reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("Failed to commit message: %v", err)
			continue
		}

		log.Printf("Message committed successfully: %s", string(msg.Value))
	}
}

func (c *Consumer) Close() {
	if err := c.Reader.Close(); err != nil {
		log.Printf("Error closing Kafka Consumer: %v", err)
	} else {
		log.Println("Kafka Consumer closed.")
	}
}
