package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"log"
	"payment_system/internal/pkg/app"
	"payment_system/internal/pkg/config"
	"payment_system/pkg/api/v1"
)

type PaymentServer struct {
	echo           *echo.Echo
	PORT           string
	paymentService *app.PaymentService
}

func NewPaymentServer() *PaymentServer {
	return &PaymentServer{}
}

func (s *PaymentServer) Init(cfg config.Config, echo *echo.Echo, PORT string) error {
	s.paymentService = app.NewPaymentService()
	err := s.paymentService.Init(cfg.System)
	if err != nil {
		log.Printf("Error initializing payment service: %v", err)
		return err
	}
	s.echo = echo
	s.PORT = PORT
	api.RegisterHandlers(s.echo, s)
	return nil
}

func (s *PaymentServer) Run() error {
	err := s.echo.Start(s.PORT)
	return err
}

func (s *PaymentServer) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
