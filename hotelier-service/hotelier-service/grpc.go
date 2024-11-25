package main

import (
	"google.golang.org/grpc"
	"hotelier-service/internal/pkg/api"
	"hotelier-service/internal/pkg/proto"
	"hotelier-service/internal/pkg/repository"
	_ "log"
	"net"
	_ "net/url"
)

func startGRPCServer(service *repository.Service) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", ":50051") // Порт gRPC сервера
	if err != nil {
		return nil, nil, err
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
	return grpcServer, lis, err
}
