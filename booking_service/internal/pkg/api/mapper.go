package api

import (
	"booking_service/pkg/api/v1"
	"booking_service/pkg/models"
	"fmt"
	"log"
	"time"
)

func fromApiBookingRequest(bookingRequest *api.BookingRequest) (*models.BookingRequest, error) {
	formattedCheckInDate, err := time.Parse("2006-01-02", bookingRequest.CheckInDate)
	if err != nil {
		log.Printf("Failed to parse checkin date: %v", err)
		return nil, err
	}

	formattedCheckOutDate, err := time.Parse("2006-01-02", bookingRequest.CheckOutDate)
	if err != nil {
		log.Printf("Failed to parse checkout date: %v", err)
		return nil, err
	}
	if formattedCheckOutDate.Before(formattedCheckInDate) {
		return nil, fmt.Errorf("checkin date %v cannot be after checkout date %v", formattedCheckInDate, formattedCheckOutDate)
	}
	return &models.BookingRequest{
		CheckInDate:         formattedCheckInDate,
		ClientFullName:      bookingRequest.ClientFullName,
		ClientPhoneNumber:   bookingRequest.ClientPhoneNumber,
		HotelierPhoneNumber: bookingRequest.HotelierPhoneNumber,
		CheckOutDate:        formattedCheckOutDate,
		HotelId:             bookingRequest.HotelId,
		RoomNumber:          bookingRequest.RoomNumber,
	}, nil
}

func toApiBooking(booking *models.Booking) *api.Booking {
	formattedCheckInDate := booking.CheckInDate.Format("2006-01-02")
	formattedCheckOutDate := booking.CheckInDate.Format("2006-01-02")

	return &api.Booking{
		BookingId:           &booking.BookingId,
		CheckInDate:         &formattedCheckInDate,
		ClientFullName:      &booking.ClientFullName,
		ClientPhoneNumber:   &booking.ClientPhoneNumber,
		HotelierPhoneNumber: &booking.HotelierPhoneNumber,
		CheckOutDate:        &formattedCheckOutDate,
		HotelId:             &booking.HotelId,
		RoomNumber:          &booking.RoomNumber,
		TotalPrice:          &booking.TotalPrice,
		PaymentStatus:       (*api.BookingPaymentStatus)(&booking.PaymentStatus),
	}
}
