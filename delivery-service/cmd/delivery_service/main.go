package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	my_grpc "delivery-service/internal/grpc"
	pb "delivery-service/internal/grpc/proto"
	"delivery-service/internal/sms"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		log.Fatalf("GRPC_PORT environment variable is not set")
	}

	// Создаем listener
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Инициализируем SMSClient
	smsClient, err := sms.NewSMSClient()
	if err != nil {
		log.Fatalf("Failed to initialize SMSClient: %v", err)
	}

	// Инициализируем DeliveryHandler
	deliveryHandler := &my_grpc.DeliveryHandler{
		SMSClient: smsClient,
	}

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(grpcServer, deliveryHandler)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		log.Printf("Delivery Service is running on port %s", port)
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	sig := <-sigs
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Завершаем gRPC сервер
	shutdownGrpcServer(grpcServer)

	<-ctx.Done()

	log.Println("Shutdown complete.")
}

func shutdownGrpcServer(grpcServer *grpc.Server) {
	log.Println("Shutting down gRPC server...")
	go grpcServer.GracefulStop()
	log.Println("gRPC server stopped gracefully.")
}
