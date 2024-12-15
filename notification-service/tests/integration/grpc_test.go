package integration

//
//import (
//	"context"
//	"testing"
//	"time"
//	pb "notification-service/internal/grpc/proto"
//	"github.com/stretchr/testify/assert"
//	"google.golang.org/grpc"
//)
//
//func TestGRPCNotification(t *testing.T) {
//	conn, err := grpc.Dial("notification-service:50051", grpc.WithInsecure()) // Подключение к gRPC-серверу
//	assert.NoError(t, err)
//	defer conn.Close()
//
//	client := pb.NewDeliveryServiceClient(conn)
//
//	request := &pb.NotificationRequest{
//		ClientPhone:   "1234567890",
//		HotelierPhone: "0987654321",
//		CheckInDate:   "2024-01-01",
//		CheckOutDate:  "2024-01-10",
//		RoomNumber:    "101",
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	response, err := client.SendNotification(ctx, request)
//	assert.NoError(t, err)                                             // Проверяем, что вызов прошел успешно
//	assert.Equal(t, "Notification sent successfully", response.Status) // Проверяем ответ
//}
