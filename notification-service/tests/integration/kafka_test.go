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
//	produceConsumeTopic = "produce-consume-topic"
//)
//
//func TestKafkaProduceConsume(t *testing.T) {
//	log.Println("Waiting for Kafka to initialize...")
//	time.Sleep(3 * time.Second)
//
//	messageKey := "test-key"
//	messageValue := map[string]interface{}{"test_field": "test_value"}
//	messageJSON, err := json.Marshal(messageValue)
//	assert.NoError(t, err)
//
//	writer := kafka.NewWriter(kafka.WriterConfig{
//		Brokers: []string{"kafka:9092"},
//		Topic:   produceConsumeTopic,
//	})
//	defer writer.Close()
//
//	err = writer.WriteMessages(context.Background(), kafka.Message{
//		Key:   []byte(messageKey),
//		Value: messageJSON,
//	})
//	assert.NoError(t, err)
//
//	reader := kafka.NewReader(kafka.ReaderConfig{
//		Brokers:  []string{"kafka:9092"},
//		Topic:    produceConsumeTopic,
//		GroupID:  "produce-consume-group",
//		MinBytes: 10e3,
//		MaxBytes: 10e6,
//	})
//	defer reader.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
//	defer cancel()
//
//	msg, err := reader.ReadMessage(ctx)
//	assert.NoError(t, err)
//	assert.Equal(t, messageKey, string(msg.Key))
//	assert.Equal(t, string(messageJSON), string(msg.Value))
//}
