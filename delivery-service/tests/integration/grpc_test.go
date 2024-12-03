package integration

import (
	"context"
	"net"
	"testing"

	grpcinternal "delivery-service/internal/grpc"
	pb "delivery-service/internal/grpc/proto"

	"google.golang.org/grpc"
)

func TestDeliveryHandler_SendNotification(t *testing.T) {
	// Создаём слушатель
	lis, err := net.Listen("tcp", ":0") // Используем случайный порт
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	// Создаём gRPC сервер
	s := grpc.NewServer()
	pb.RegisterDeliveryServiceServer(s, &grpcinternal.DeliveryHandler{})
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatalf("failed to serve: %v", err)
		}
	}()
	defer s.Stop()

	// Создаём клиента
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDeliveryServiceClient(conn)

	// Тестируем SendNotification
	req := &pb.NotificationRequest{
		ClientPhone:   "1234567890",
		HotelierPhone: "0987654321",
		CheckInDate:   "2024-12-01",
		CheckOutDate:  "2024-12-10",
		RoomNumber:    "101",
	}

	resp, err := client.SendNotification(context.Background(), req)
	if err != nil {
		t.Fatalf("SendNotification failed: %v", err)
	}

	if resp.Status != "Notification sent successfully" {
		t.Errorf("unexpected response: %v", resp.Status)
	}
}
