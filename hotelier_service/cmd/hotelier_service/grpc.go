package main

import (
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/app"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/pkg/api"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/proto"
	"google.golang.org/grpc"
	_ "log"
	"net"
	_ "net/url"
)

func startGRPCServer(service app.HotelService, config Config) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+config.Server.InternalGrpcPort) // Порт gRPC сервера
	if err != nil {
		return nil, nil, err
	}
	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()
	hotelierServer := api.NewHotelierServer(service)
	// Регистрируем реализацию сервиса
	proto.RegisterHotelierServiceServer(grpcServer, hotelierServer)

	//log.Println("gRPC server started on :50051", lis.Addr())
	//if err := grpcServer.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	return grpcServer, lis, err
}
