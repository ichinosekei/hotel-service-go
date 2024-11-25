package app

import (
	"booking_service/internal/app"
	"booking_service/internal/pkg/config"
	"booking_service/internal/pkg/persistent/repository"
	"booking_service/internal/pkg/persistent/server"
	"context"
	"log"
)

type BookingService struct {
	service *app.Service
}

func NewBookingService() *BookingService {
	return &BookingService{}
}
func (s *BookingService) Init(cfg config.Config) error {
	repo := repository.NewRepository()
	err := repo.Init(cfg.Database)
	if err != nil {
		log.Fatal("Failed to init repository")
		return err
	}

	srv := server.NewBookingServer(repo)
	srv.Init(cfg.Server)

	s.service = app.NewService(srv, repo)
	return nil
}
func (s *BookingService) Run() error {
	return s.service.Run()
}

func (s *BookingService) Shutdown(ctx context.Context) error {
	return s.service.Shutdown(ctx)
}
