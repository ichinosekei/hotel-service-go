package grpc

import (
	"context"
	"errors"
	"log"
	"notification-service/internal/grpc/proto"
	"time"
)

func SendNotification(client proto.DeliveryServiceClient, ctx context.Context, notification *proto.NotificationRequest) (*proto.NotificationResponse, error) {
	// Таймаут 10 секунд для gRPC вызова
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := client.SendNotification(ctx, notification)
	if err != nil {
		log.Printf("Error sending gRPC notification: %v", err)
		return nil, errors.New("unexpected nil") // Здесь была ошибка: порядок возврата был неверным
	}
	log.Printf("Notification sent successfully: %s", resp.Status)
	return resp, nil
}
