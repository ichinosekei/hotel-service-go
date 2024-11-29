package app

import (
	"booking_service/internal/app"
	"booking_service/internal/pkg/persistent/repository"
	"log"
)

type BookingService struct {
	Service *app.Service
}

func NewBookingService() *BookingService {
	return &BookingService{}
}
func (s *BookingService) Init(cfg repository.Config) error {
	repo := repository.NewRepository()
	err := repo.Init(cfg)
	if err != nil {
		log.Fatal("Failed to init repository")
		return err
	}
	s.Service = app.NewService(repo)
	return nil
}
