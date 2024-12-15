package kafka

import (
	"github.com/segmentio/kafka-go"
)

type BookingEvent struct {
	ClientPhone string `json:"client_phone"`
	HotelPhone  string `json:"hotel_phone"`
	CheckIn     string `json:"check_in"`
	CheckOut    string `json:"check_out"`
	Room        string `json:"room"`
	EventID     string `json:"event_id"`
}

func NewConsumer(brokers []string, topic string, groupID string) (*Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})

	return &Consumer{Reader: reader}, nil
}

//func NewProducer(brokers, topic string) (kafka.Writer, error) {
//	return kafka.Writer{
//		Addr:     kafka.TCP(brokers),
//		Topic:    topic,
//		Balancer: &kafka.LeastBytes{},
//	}, nil
//}
