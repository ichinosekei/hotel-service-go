package repository

import "booking_service/pkg/api/v1"

func toModelsBooking(booking *Booking) *api.Booking {
	formattedCheckInDate := booking.CheckInDate.Format("2006-01-02")
	formattedCheckOutDate := booking.CheckInDate.Format("2006-01-02")

	return &api.Booking{
		BookingId:         &booking.BookingId,
		CheckInDate:       &formattedCheckInDate,
		ClientFullName:    &booking.ClientFullName,
		ClientPhoneNumber: &booking.ClientPhoneNumber,
		CheckOutDate:      &formattedCheckOutDate,
		HotelId:           &booking.HotelId,
		RoomNumber:        &booking.RoomNumber,
		TotalPrice:        &booking.TotalPrice,
	}
}
