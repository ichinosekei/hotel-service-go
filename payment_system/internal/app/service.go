package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"payment_system/internal/pkg/config"
	"payment_system/pkg/models"
)

type PaymentService struct {
	addr string
}

func NewPaymentService(cfg config.Config) *PaymentService {
	return &PaymentService{
		addr: cfg.Addr,
	}
}

func (s *PaymentService) Create(request *models.PaymentRequest) (*models.PaymentResponse, error) {
	paymentID := uuid.NewString()
	response := models.PaymentResponse{
		URL:       "https://payment/" + paymentID,
		PaymentID: paymentID,
	}
	return &response, nil
}

func (s *PaymentService) Send(request *models.PaymentWebhookRequest) error {
	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal Payment Webhook request: %v", err)
		return err
	}
	resp, err := http.Post(s.addr, "application/json", bytes.NewReader(body))
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
