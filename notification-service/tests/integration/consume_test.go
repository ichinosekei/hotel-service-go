package integration

//
//import (
//	"context"
//	"encoding/json"
//	"log"
//	"testing"
//	"time"
//	"github.com/segmentio/kafka-go"
//	"github.com/stretchr/testify/assert"
//)
//
//const (
//	consumeTopic   = "consume-booking-events" // Unique topic for consume tests
//	consumeGroupID = "consume-test-group"     // Unique consumer group
//)
//
//// BookingEvent represents the schema for booking event messages
//
//type BookingEvent struct {
//    EventID     string `json:"event_id"`
//    ClientPhone string `json:"client_phone"`
//    CheckIn     string `json:"check_in"`
//    CheckOut    string `json:"check_out"`
//    Room        string `json:"room"`
//    HotelPhone  string `json:"hotel_phone"`
//}
//	ClientPhone string `json:"client_phone"`
//	CheckIn     string `json:"check_in"`
//	CheckOut    string `json:"check_out"`
//	Room        string `json:"room"`
//	HotelPhone  string `json:"hotel_phone"`
//}
//
//func TestKafkaConsume(t *testing.T) {
//	// Test messages to publish
//	messages := []BookingEvent{
//	EventID: "test-event-id",
//		{ClientPhone: "1234567890", CheckIn: "2024-01-01", CheckOut: "2024-01-10", Room: "101", HotelPhone: "0987654321"},
//		{ClientPhone: "9876543210", CheckIn: "2024-02-01", CheckOut: "2024-02-10", Room: "202", HotelPhone: "1234567890"},
//	}
//
//	writer := kafka.NewWriter(kafka.WriterConfig{
//		Brokers:  []string{"kafka:9092"},
//		Topic:    consumeTopic,
//		Balancer: &kafka.LeastBytes{},
//	})
//	defer writer.Close()
//
//	for _, msg := range messages {
//		messageJSON, err := json.Marshal(msg)
//		assert.NoError(t, err)
//
//		err = writer.WriteMessages(context.Background(), kafka.Message{
//			Key:   []byte(msg.ClientPhone),
//			Value: messageJSON,
//		})
//		assert.NoError(t, err)
//		log.Printf("Published message: %s", messageJSON)
//	}
//
//	reader := kafka.NewReader(kafka.ReaderConfig{
//		Brokers:  []string{"kafka:9092"},
//		Topic:    consumeTopic,
//		GroupID:  consumeGroupID,
//		MinBytes: 10e3,
//		MaxBytes: 10e6,
//	})
//	defer reader.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
//	defer cancel()
//
//	var receivedMessages []BookingEvent
//	for i := 0; i < len(messages); i++ {
//		msg, err := reader.ReadMessage(ctx)
//		assert.NoError(t, err)
//
//		var event BookingEvent
//		err = json.Unmarshal(msg.Value, &event)
//		assert.NoError(t, err)
//		receivedMessages = append(receivedMessages, event)
//	}
//
//	assert.Equal(t, messages, receivedMessages, "Published and received messages should match")
//}
