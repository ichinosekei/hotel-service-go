package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"payment_system/internal/app"
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

func (s *PaymentServer) Init(cfg config.Config, echo *echo.Echo, PORT string) {
	s.paymentService = app.NewPaymentService(cfg)
	s.echo = echo
	s.PORT = PORT
	api.RegisterHandlers(s.echo, s)
}

func (s *PaymentServer) Run() error {
	err := s.echo.Start(s.PORT)
	return err
}

func (s *PaymentServer) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
