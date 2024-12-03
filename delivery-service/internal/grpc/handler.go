package grpc

import (
	"context"
	"delivery-service/internal/sms"
	"log"

	pb "delivery-service/internal/grpc/proto"
)

// gRPC-обработчик
type DeliveryHandler struct {
	pb.UnimplementedDeliveryServiceServer
	SMSClient *sms.SMSClient
}

// Отправляет уведомление клиенту и отельеру
func (h *DeliveryHandler) SendNotification(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	clientMessage := "Ваше бронирование подтверждено. Даты: " + req.CheckInDate + " - " + req.CheckOutDate
	hotelMessage := "Поступило новое бронирование. Клиент: " + req.ClientPhone

	// SMS клиенту
	if err := h.SMSClient.SendSMS(req.ClientPhone, clientMessage); err != nil {
		log.Printf("Failed to send SMS to client: %v", err)
		return nil, err
	}

	// SMS отельеру
	if err := h.SMSClient.SendSMS(req.HotelierPhone, hotelMessage); err != nil {
		log.Printf("Failed to send SMS to hotelier: %v", err)
		return nil, err
	}

	log.Printf("Notifications sent successfully to client %s and hotelier %s", req.ClientPhone, req.HotelierPhone)
	return &pb.NotificationResponse{Status: "Notification sent successfully"}, nil
}
