package repository

import (
	"booking_service/pkg/api/v1"
	"fmt"
	"log"
	"time"
)

func LoadBookingRequest(bookingRequest *api.BookingRequest) (*Booking, error) {
	formattedCheckInDate, err := time.Parse("2006-01-02", *bookingRequest.CheckInDate)
	if err != nil {
		log.Printf("Failed to parse checkin date: %v", err)
		return &Booking{}, err
	}

	formattedCheckOutDate, err := time.Parse("2006-01-02", *bookingRequest.CheckOutDate)
	if err != nil {
		log.Printf("Failed to parse checkout date: %v", err)
		return &Booking{}, err
	}
	if formattedCheckOutDate.Before(formattedCheckInDate) {
		return &Booking{}, fmt.Errorf("checkin date %v cannot be after checkout date %v", formattedCheckInDate, formattedCheckOutDate)
	}
	return &Booking{
		CheckInDate:       formattedCheckInDate,
		ClientFullName:    *bookingRequest.ClientFullName,
		ClientPhoneNumber: *bookingRequest.ClientPhoneNumber,
		CheckOutDate:      formattedCheckOutDate,
		HotelId:           *bookingRequest.HotelId,
		RoomNumber:        *bookingRequest.RoomNumber,
	}, nil
}
