package server

import (
	"booking_service/internal/app"
	"booking_service/internal/app/server_gen"
	"github.com/labstack/echo/v4"
)

type BookingServer struct {
	repo app.Repository
	echo *echo.Echo
	URL  string
}

func NewBookingServer(repo app.Repository) *BookingServer {
	return &BookingServer{repo, nil, ""}
}

func (s *BookingServer) Init(cfg Config) {
	s.echo = echo.New()
	s.URL = cfg.URL
	server_gen.RegisterHandlers(s.echo, s)
}

func (s *BookingServer) Run() error {
	err := s.echo.Start(s.URL)
	return err
}
