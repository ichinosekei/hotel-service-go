package main

import (
	"log"
	"net"
	_ "net/url"
     "hotel-service-go/hotelier-service/internal/pkg/repository"
	"google.golang.org/grpc"
	"hotel-service-go/hotelier-service/internal/pkg/api"
	"hotel-service-go/hotelier-service/internal/pkg/proto"
)

func startGRPCServer(service *repository.Service) (*grpc.Server, net.Listener){
	lis, err := net.Listen("tcp", ":50051")// Порт gRPC сервера
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()
	hotelierServer := api.NewHotelierServer(service)
	// Регистрируем реализацию сервиса
	proto.RegisterHotelierServiceServer(grpcServer, hotelierServer)

	//log.Println("gRPC server started on :50051")
	//if err := grpcServer.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	return grpcServer, lis
}

