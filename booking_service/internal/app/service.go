package app

import (
	"booking_service/pkg/api/v1"
	"booking_service/pkg/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetClient(phoneNumber string) (*models.Bookings, error) {
	return s.repo.GetClient(phoneNumber)
}

func (s *Service) CreateClient(bookingRequest *api.BookingRequest) error {
	return s.repo.Create(bookingRequest)
}

func (s *Service) GetHotel(hotelId int) (*models.Bookings, error) {
	return s.repo.GetHotel(hotelId)
}
