package app

import (
	"booking_service/pkg/models"
)

type Repository interface {
	Create(*models.BookingRequest) error
	GetClient(string) (*models.Bookings, error)
	GetHotel(int) (*models.Bookings, error)
}
