package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"payment_system/pkg/models"
)

type PaymentSystem struct {
	bookingAddr string
}

func NewPaymentSystem(cfg Config) *PaymentSystem {
	return &PaymentSystem{
		bookingAddr: cfg.BookingAddr,
	}
}

func (s *PaymentSystem) Create(paymentRequest *models.PaymentRequest) (*models.PaymentResponse, error) {
	paymentID := uuid.NewString()
	response := PaymentResponse{
		"0",
		paymentID,
	}
	return toModelPaymentResponse(&response), nil
}

func (s *PaymentSystem) Send(paymentRequest *models.PaymentWebhookRequest) error {
	request := fromModelPaymentWebhookRequest(paymentRequest)
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal Payment Webhook request: %v", err)
		return err
	}
	resp, err := http.Post(s.bookingAddr, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("Failed to book Payment Webhook request: %v", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to book Payment Webhook request: %v", resp.StatusCode)
		return fmt.Errorf("failed to book Payment Webhook request: %v", resp.StatusCode)
	}
	return nil
}
