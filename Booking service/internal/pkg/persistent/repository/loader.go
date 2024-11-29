package repository

import (
	"booking_service/pkg/models"
)

func LoadBookingRequest(bookingRequest *models.BookingRequest) *Booking {
	return &Booking{
		CheckInDate:       bookingRequest.CheckInDate,
		ClientFullName:    bookingRequest.ClientFullName,
		ClientPhoneNumber: bookingRequest.ClientPhoneNumber,
		CheckOutDate:      bookingRequest.CheckOutDate,
		HotelId:           bookingRequest.HotelId,
		RoomNumber:        bookingRequest.RoomNumber,
	}
}
