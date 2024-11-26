package repository

import (
	"booking_service/pkg/api/v1"
)

func toModelsBooking(booking *Booking) *api.Booking {
	return &api.Booking{
		BookingId:         &booking.BookingId,
		CheckInDate:       &booking.CheckInDate,
		ClientFullName:    &booking.ClientFullName,
		ClientPhoneNumber: &booking.ClientPhoneNumber,
		Duration:          &booking.Duration,
		HotelId:           &booking.HotelId,
		RoomNumber:        &booking.RoomNumber,
		TotalPrice:        &booking.TotalPrice,
	}
}
