package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "notification-service/internal/grpc/proto"
)

type MockDeliveryService struct {
	pb.UnimplementedDeliveryServiceServer
}

func (m *MockDeliveryService) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	log.Printf("Mock DeliveryService: Sending notification to client %s and hotelier %s", req.ClientPhone, req.HotelierPhone)
	return &pb.NotificationResponse{
		Status: "Notification sent successfully",
	}, nil
}

func StartMockDeliveryService(port string) error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(server, &MockDeliveryService{})

	log.Printf("Mock DeliveryService running on %s", port)
	return server.Serve(listener)
}
