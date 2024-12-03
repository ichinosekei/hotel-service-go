package main

import (
	"log"
	"net"
	"os"

	my_grpc "delivery-service/internal/grpc"
	pb "delivery-service/internal/grpc/proto"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	// Настройка gRPC-сервера
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(grpcServer, &my_grpc.DeliveryHandler{})

	log.Printf("Delivery Service is running on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
