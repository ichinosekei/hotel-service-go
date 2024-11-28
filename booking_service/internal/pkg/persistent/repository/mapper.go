package repository

import (
	"booking_service/pkg/models"
)

func toModelsBooking(booking *Booking) *models.Booking {
	return &models.Booking{
		BookingId:         booking.BookingId,
		CheckInDate:       booking.CheckInDate,
		ClientFullName:    booking.ClientFullName,
		ClientPhoneNumber: booking.ClientPhoneNumber,
		CheckOutDate:      booking.CheckOutDate,
		HotelId:           booking.HotelId,
		RoomNumber:        booking.RoomNumber,
		TotalPrice:        booking.TotalPrice,
	}
}
