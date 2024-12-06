package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	pb "notification-service/internal/grpc/proto"
	"notification-service/internal/kafka"
)

func main() {
	log.Println("Notification Service starting...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received signal: %v. Initiating shutdown...", sig)
		cancel()
	}()

	kafkaBrokers, kafkaTopic, grpcAddress := readEnvVars()

	// готовность Kafka
	if err := waitForKafkaTopic(kafkaBrokers, kafkaTopic, 60, 3*time.Second); err != nil {
		log.Fatalf("Kafka topic '%s' is not ready: %v", kafkaTopic, err)
	}

	consumer, err := kafka.NewConsumer(kafkaBrokers, kafkaTopic, "notification-service-group")
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}

	grpcClient, conn := initializeGrpcClient(grpcAddress)

	deduplicатор := kafka.NewMessageDeduplicator()

	// Запуск Kafka Consumer
	go consumer.Start(ctx, func(event kafka.BookingEvent) error {
		return processBookingEvent(ctx, grpcClient, deduplicатор, event)
	})

	// Завершаем подребителя
	go shutdownKafka(ctx, consumer)

	<-ctx.Done()

	// Завершаем ресурсы
	log.Println("Shutting down gracefully...")
	shutdownGrpcClient(conn)
	log.Println("Shutdown complete.")
}

func readEnvVars() ([]string, string, string) {
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

	return []string{kafkaBrokers}, kafkaTopic, grpcAddress
}

func waitForKafkaTopic(brokers []string, topic string, retries int, retryInterval time.Duration) error {
	log.Printf("Waiting for Kafka topic '%s' to become available...", topic)

	consumer, err := kafka.NewConsumer(brokers, topic, "health-check")
	if err != nil {
		return fmt.Errorf("failed to initialize Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Проверяем доступность топика за ограниченное число попыток
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), retryInterval)
		defer cancel()

		// Пробуем читать сообщение
		_, err := consumer.Reader.FetchMessage(ctx)
		if err == nil {
			log.Printf("Kafka topic '%s' is ready.", topic)
			return nil
		}

		log.Printf("Attempt %d/%d: Kafka topic '%s' not ready. Error: %v", i+1, retries, topic, err)
		time.Sleep(retryInterval)
	}

	return fmt.Errorf("Kafka topic '%s' is not ready after %d attempts", topic, retries)
}

func initializeGrpcClient(address string) (pb.DeliveryServiceClient, *grpc.ClientConn) {
	log.Println("Connecting to gRPC delivery service...")
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to delivery service: %v", err)
	}
	log.Println("gRPC client connected successfully.")
	return pb.NewDeliveryServiceClient(conn), conn
}

func shutdownGrpcClient(conn *grpc.ClientConn) {
	log.Println("Shutting down gRPC client...")
	if err := conn.Close(); err != nil {
		log.Printf("Error closing gRPC connection: %v", err)
	} else {
		log.Println("gRPC connection closed.")
	}
}

// Завершение потребителя через select case
func shutdownKafka(ctx context.Context, consumer *kafka.Consumer) {
	log.Println("Waiting to shutdown Kafka Consumer...")
	select {
	case <-ctx.Done():
		log.Println("Context canceled, closing Kafka Consumer...")
		consumer.Close() // Вызываем Close() без ожидания возвращаемого значения
		log.Println("Kafka Consumer closed successfully.")
	}
}

// обработка событий BookingEvent
func processBookingEvent(ctx context.Context, grpcClient pb.DeliveryServiceClient, deduplicатор *kafka.MessageDeduplicator, event kafka.BookingEvent) error {
	if deduplicатор.IsDuplicate(event.EventID) {
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
}
