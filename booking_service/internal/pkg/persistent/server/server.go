package server

import (
	"booking_service/internal/app"
	"booking_service/internal/app/server_gen"
	"context"
	"github.com/labstack/echo/v4"
)

type BookingServer struct {
	repo app.Repository
	echo *echo.Echo
	PORT string
}

func NewBookingServer(repo app.Repository) *BookingServer {
	return &BookingServer{repo, echo.New(), ""}
}

func (s *BookingServer) Init(cfg Config) {
	s.PORT = cfg.PORT
	server_gen.RegisterHandlers(s.echo, s)
}

func (s *BookingServer) Run() error {
	err := s.echo.Start(s.PORT)
	return err
}

func (s *BookingServer) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(context.Background())
}
