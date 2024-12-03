package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	pb "notification-service/internal/grpc/proto"
	"notification-service/internal/kafka"
)

func main() {
	log.Println("Notification Service starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandler(cancel)

	kafkaBrokers, kafkaTopic, grpcAddress := readEnvVars()

	// Проверяем готовность Kafka (чтоб не писать в неоткрытые топики, создавать топик при загрузке не получилось)
	waitForKafkaTopic(kafkaBrokers, kafkaTopic)

	// Kafka Consumer
	consumer := initializeKafkaConsumer(kafkaBrokers, kafkaTopic)
	defer consumer.Reader.Close()

	grpcClient := initializeGrpcClient(grpcAddress)

	// Чтоб не дублировать сообщения
	deduplicator := kafka.NewMessageDeduplicator()

	// Kafka Consumer
	startConsumer(ctx, consumer, grpcClient, deduplicator)

	<-ctx.Done()
	log.Println("Shutting down...")
}

// Обработчик сигналов для корректного завершения приложения.
func setupSignalHandler(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()
}

func readEnvVars() (string, string, string) {
	kafkaBrokers := os.Getenv("KAFKA_BROKER")
	if kafkaBrokers == "" {
		log.Fatalf("KAFKA_BROKER is not set")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		log.Fatalf("KAFKA_TOPIC is not set")
	}

	grpcAddress := os.Getenv("DELIVERY_SERVICE_ADDRESS")
	if grpcAddress == "" {
		log.Fatalf("DELIVERY_SERVICE_ADDRESS is not set")
	}

	return kafkaBrokers, kafkaTopic, grpcAddress
}

func initializeKafkaConsumer(brokers, topic string) *kafka.Consumer {
	log.Println("Initializing Kafka Consumer...")
	consumer, err := kafka.NewConsumer([]string{brokers}, topic, "notification-service-group")
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}
	log.Println("Kafka Consumer initialized.")
	return consumer
}

func initializeGrpcClient(address string) pb.DeliveryServiceClient {
	log.Println("Connecting to gRPC delivery service...")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to delivery service: %v", err)
	}
	log.Println("gRPC client connected successfully.")
	return pb.NewDeliveryServiceClient(conn)
}

func waitForKafkaTopic(broker, topic string) {
	log.Printf("Waiting for Kafka topic '%s' to become available...", topic)
	config := sarama.NewConfig()
	config.Metadata.Retry.Backoff = 500 * time.Millisecond
	config.Metadata.RefreshFrequency = 2 * time.Second

	for {
		client, err := sarama.NewClient([]string{broker}, config)
		if err != nil {
			log.Printf("[waitForKafkaTopic] Kafka unavailable: %v. Retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer client.Close()

		topics, err := client.Topics()
		if err != nil {
			log.Printf("[waitForKafkaTopic] Error fetching topics: %v. Retrying...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if topicExists(topics, topic) {
			log.Printf("[waitForKafkaTopic] Topic '%s' is ready.", topic)
			return
		}
		log.Printf("[waitForKafkaTopic] Topic '%s' not found. Retrying...", topic)
		time.Sleep(5 * time.Second)
	}
}

// Проверяет, существует ли топик в списке.
func topicExists(topics []string, topic string) bool {
	for _, t := range topics {
		if t == topic {
			return true
		}
	}
	return false
}

func startConsumer(ctx context.Context, consumer *kafka.Consumer, grpcClient pb.DeliveryServiceClient, deduplicator *kafka.MessageDeduplicator) {
	log.Println("Starting Kafka Consumer...")
	go consumer.StartConsumer(ctx, func(event kafka.BookingEvent) error {
		if deduplicator.IsDuplicate(event.EventID) {
			log.Printf("Duplicate event detected: %s", event.EventID)
			return nil
		}

		notification := &pb.NotificationRequest{
			ClientPhone:   event.ClientPhone,
			HotelierPhone: event.HotelPhone,
			CheckInDate:   event.CheckIn,
			CheckOutDate:  event.CheckOut,
			RoomNumber:    event.Room,
		}

		resp, err := grpcClient.SendNotification(ctx, notification)
		if err != nil {
			log.Printf("Failed to send notification: %v", err)
			return err
		}
		log.Printf("Notification sent successfully: %s", resp.Status)
		return nil
	})
}
