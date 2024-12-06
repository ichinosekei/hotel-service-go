package api

import (
	"booking_service/pkg/api/v1"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func (s *BookingServer) PostApiV1Bookings(ctx echo.Context) error {
	var bookingRequest api.BookingRequest

	if err := ctx.Bind(&bookingRequest); err != nil {
		log.Printf("Failed to binding bookings request: %v\n", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}
	modelRequest, err := fromApiBookingRequest(&bookingRequest)
	if err != nil {
		log.Printf("Failed to process bookings request: %v\n", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}
	if err = s.bookingService.Service.CreateClient(modelRequest); err != nil {
		log.Printf("Error creating booking: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, &bookingRequest)
}

func (s *BookingServer) GetApiV1BookingsClient(ctx echo.Context, params api.GetApiV1BookingsClientParams) error {
	bookings, err := s.bookingService.Service.GetClient(params.PhoneNumber)
	if err != nil {
		log.Printf("Error getting bookings: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	var apiBookings []api.Booking
	for _, booking := range *bookings {
		apiBookings = append(apiBookings, *toApiBooking(&booking))
	}

	return ctx.JSON(http.StatusOK, apiBookings)
}

func (s *BookingServer) GetApiV1BookingsHotel(ctx echo.Context, params api.GetApiV1BookingsHotelParams) error {
	bookings, err := s.bookingService.Service.GetHotel(params.HotelId)
	if err != nil {
		log.Printf("Error getting bookings: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	var apiBookings []api.Booking
	for _, booking := range *bookings {
		apiBookings = append(apiBookings, *toApiBooking(&booking))
	}

	return ctx.JSON(http.StatusOK, apiBookings)
}

func (s *BookingServer) PostApiV1WebhookPayment(ctx echo.Context) error {
	var paymentRequest api.PaymentWebhookRequest
	if err := ctx.Bind(&paymentRequest); err != nil {
		log.Printf("Invalid webhook request: %v", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if paymentRequest.Status != api.Ok {
		log.Printf("Invalid status: %v", paymentRequest.Status)
		return ctx.JSON(http.StatusBadRequest, fmt.Errorf("booking has not been paid"))
	}

	err := s.bookingService.Service.UpdatePaymentStatus(paymentRequest.BookingId)
	if err != nil {
		log.Printf("Error updating booking: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	log.Printf("Updated booking: %v\n", paymentRequest.BookingId)
	return ctx.JSON(http.StatusOK, nil)
}
