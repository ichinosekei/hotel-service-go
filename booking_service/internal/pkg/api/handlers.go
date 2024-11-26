package api

import (
	"booking_service/pkg/api/v1"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func (s *BookingServer) PostApiV1Bookings(ctx echo.Context) error {
	var bookingRequest api.BookingRequest

	if err := ctx.Bind(&bookingRequest); err != nil {
		log.Printf("Error binding bookings request: %v\n", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	if err := s.bookingService.Service.CreateClient(&bookingRequest); err != nil {
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
	return ctx.JSON(http.StatusOK, bookings)
}

func (s *BookingServer) GetApiV1BookingsHotel(ctx echo.Context, params api.GetApiV1BookingsHotelParams) error {
	bookings, err := s.bookingService.Service.GetHotel(params.HotelId)
	if err != nil {
		log.Printf("Error getting bookings: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, bookings)
}
