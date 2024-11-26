package app

import (
	"booking_service/pkg/api/v1"
	"booking_service/pkg/models"
)

type Repository interface {
	Create(*api.BookingRequest) error
	GetClient(string) (*models.Bookings, error)
	GetHotel(int) (*models.Bookings, error)
}
