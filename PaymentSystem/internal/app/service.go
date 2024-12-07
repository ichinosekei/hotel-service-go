package app

import "payment_system/pkg/models"

type Service struct {
	sys System
}

func NewService(sys System) *Service {
	return &Service{sys: sys}
}

func (s *Service) Create(paymentRequest *models.PaymentRequest) (*models.PaymentResponse, error) {
	return s.sys.Create(paymentRequest)
}

func (s *Service) Send(paymentRequest *models.PaymentWebhookRequest) error {
	return s.sys.Send(paymentRequest)
}
