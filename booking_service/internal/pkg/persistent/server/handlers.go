package server

import (
	"booking_service/internal/app/server_gen"
	"booking_service/internal/pkg/persistent/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func (s *BookingServer) PostApiV1Bookings(ctx echo.Context) error {
	var request server_gen.BookingRequest

	if err := ctx.Bind(&request); err != nil {
		log.Printf("Error binding bookings request: %v\n", err)
		return ctx.JSON(http.StatusBadRequest, err)
	}

	booking := repository.LoadBookingRequest(request)
	booking.ID = uuid.NewString()
	// TODO implement a grpc request to the hotel service
	log.Printf("Room price: %f")
	booking.TotalPrice, err := hotelierClient.GetRoomPrice(1, "101")
	if err != nil {
		log.Fatalf("Failed to get room price: %v", err)
	}
	//log.Printf("Room price: %f", price)

	if err := s.repo.Create(booking); err != nil {
		log.Printf("Error creating booking: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, &booking)
}

func (s *BookingServer) GetApiV1BookingsClient(ctx echo.Context, params server_gen.GetApiV1BookingsClientParams) error {
	bookings, err := s.repo.GetClient(params.PhoneNumber)
	if err != nil {
		log.Printf("Error getting bookings: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, bookings)
}

func (s *BookingServer) GetApiV1BookingsHotel(ctx echo.Context, params server_gen.GetApiV1BookingsHotelParams) error {
	bookings, err := s.repo.GetHotel(params.HotelId)
	if err != nil {
		log.Printf("Error getting bookings: %v\n", err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, bookings)
}
