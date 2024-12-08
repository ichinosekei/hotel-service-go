package app

import (
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

func (s *Service) GetHotel(hotelId int) (*models.Bookings, error) {
	return s.repo.GetHotel(hotelId)
}

func (s *Service) CreateClient(bookingRequest *models.BookingRequest) error {
	return s.repo.Create(bookingRequest)
}

func (s *Service) UpdatePaymentStatusPaid(bookingId string) error {
	return s.repo.UpdatePaymentStatusPaid(bookingId)
}
