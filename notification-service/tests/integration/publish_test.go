package integration

//
//import (
//	"context"
//	"encoding/json"
//	"log"
//	"testing"
//	"github.com/segmentio/kafka-go"
//	"github.com/stretchr/testify/assert"
//)
//
//const (
//	publishTopic = "publish-test-topic"
//)
//
//func TestPublishMessages(t *testing.T) {
//	messages := []map[string]interface{}{
//		{"client_phone": "1234567890", "check_in": "2024-01-01", "check_out": "2024-01-10", "room": "101", "hotel_phone": "0987654321"},
//		{"client_phone": "9876543210", "check_in": "2024-02-01", "check_out": "2024-02-10", "room": "202", "hotel_phone": "1234567890"},
//	}
//
//	writer := kafka.NewWriter(kafka.WriterConfig{
//		Brokers:  []string{"kafka:9092"},
//		Topic:    publishTopic,
//		Balancer: &kafka.LeastBytes{},
//	})
//	defer writer.Close()
//
//	for _, msg := range messages {
//		messageJSON, err := json.Marshal(msg)
//		assert.NoError(t, err)
//
//		err = writer.WriteMessages(context.Background(), kafka.Message{
//			Key:   []byte(msg["client_phone"].(string)),
//			Value: messageJSON,
//		})
//		assert.NoError(t, err)
//		log.Printf("Published message: %s", messageJSON)
//	}
//
//	log.Println("All test messages published successfully")
//}
