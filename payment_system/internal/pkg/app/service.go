package app

import (
	"payment_system/internal/app"
	"payment_system/internal/pkg/persistent/system"
)

type PaymentService struct {
	Service *app.Service
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (s *PaymentService) Init(cfg system.Config) error {
	sys := system.NewPaymentSystem(cfg)
	s.Service = app.NewService(sys)
	return nil
}
