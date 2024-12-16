package clients

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	})

	return &Producer{writer: writer}, nil
}

func (p *Producer) Send(ctx context.Context, message kafka.Message) error {
	if err := p.writer.WriteMessages(ctx, message); err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}
	log.Printf("Message sent: %v", message)
	return nil
}
func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		log.Printf("Failed to closing writer: %v", err)
		return err
	}
	log.Printf("Producer closed")
	return nil
}
