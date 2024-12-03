package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
)

func SendMessage(ctx context.Context, writer *kafka.Writer, event BookingEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.EventID),
		Value: message,
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	return errors.New("unexpected nil")
}
