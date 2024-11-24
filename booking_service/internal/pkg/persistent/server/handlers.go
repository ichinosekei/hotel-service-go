package server

import (
	"booking_service/internal/app/server_gen"
	"booking_service/internal/pkg/persistent/repository"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *BookingServer) PostApiV1Bookings(ctx echo.Context) error {
	fmt.Println("PostApiV1Bookings called")
	var request server_gen.BookingRequest

	if err := ctx.Bind(&request); err != nil {
		fmt.Println("Error parsing request:", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}
	fmt.Println("Parsed request:")
	booking := repository.LoadBookingRequest(request)
	booking.ID = 100
	booking.TotalPrice = 0
	if err := s.repo.Create(booking); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, &booking)
}

func (s *BookingServer) GetApiV1BookingsClient(ctx echo.Context, params server_gen.GetApiV1BookingsClientParams) error {
	bookings, err := s.repo.GetClient(params.PhoneNumber)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, bookings)
}

func (s *BookingServer) GetApiV1BookingsHotel(ctx echo.Context, params server_gen.GetApiV1BookingsHotelParams) error {
	bookings, err := s.repo.GetHotel(params.HotelId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, bookings)
}
