package repository

import (
	"booking_service/internal/app/server_gen"
)

func LoadBookingRequest(req server_gen.BookingRequest) *Booking {
	return &Booking{
		CheckInDate:       req.CheckInDate.Time,
		ClientFullName:    *req.ClientFullName,
		ClientPhoneNumber: *req.ClientPhoneNumber,
		Duration:          *req.Duration,
		HotelId:           *req.HotelId,
		RoomNumber:        *req.RoomNumber,
	}
}
