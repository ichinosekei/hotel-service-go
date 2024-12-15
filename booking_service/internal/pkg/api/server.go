package api

import (
	"booking_service/internal/pkg/app"
	"booking_service/internal/pkg/config"
	"booking_service/pkg/api/v1"
	"context"
	"github.com/labstack/echo/v4"
	"log"
)

type BookingServer struct {
	bookingService *app.BookingService
	echo           *echo.Echo
	PORT           string
}

func NewBookingServer() *BookingServer {
	return &BookingServer{}
}

func (s *BookingServer) Init(cfg config.Config, echo *echo.Echo, PORT string) error {
	s.echo = echo
	s.PORT = PORT
	s.bookingService = app.NewBookingService()
	err := s.bookingService.Init(cfg.Database)
	if err != nil {
		log.Fatal("Failed to init booking server")
		return err
	}
	api.RegisterHandlers(s.echo, s)
	return nil
}

func (s *BookingServer) Run() error {
	err := s.echo.Start(s.PORT)
	return err
}

func (s *BookingServer) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
