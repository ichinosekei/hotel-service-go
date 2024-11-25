package app

import (
	"booking_service/internal/pkg/persistent/repository"
)

type Repository interface {
	Create(*repository.Booking) error
	GetClient(string) (repository.Bookings, error)
	GetHotel(int) (repository.Bookings, error)
}
