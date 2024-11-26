package repository

import (
	"booking_service/pkg/api/v1"
)

func LoadBookingRequest(req *api.BookingRequest) *Booking {
	return &Booking{
		CheckInDate:       *req.CheckInDate,
		ClientFullName:    *req.ClientFullName,
		ClientPhoneNumber: *req.ClientPhoneNumber,
		Duration:          *req.Duration,
		HotelId:           *req.HotelId,
		RoomNumber:        *req.RoomNumber,
	}
}
