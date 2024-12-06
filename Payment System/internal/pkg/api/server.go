package api

import (
	"context"
	"github.com/labstack/echo/v4"
	"payment_system/pkg/api/v1"
)

type PaymentServer struct {
	echo *echo.Echo
	PORT string
}

func NewBookingServer() *PaymentServer {
	return &PaymentServer{}
}

func (s *PaymentServer) Init(echo *echo.Echo, PORT string) error {
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
