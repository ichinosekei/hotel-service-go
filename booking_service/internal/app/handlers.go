package app

import (
	openapitypes "github.com/oapi-codegen/runtime/types"
	"net/http"

	"booking_service/internal/db"
	"github.com/labstack/echo/v4"
)

type ServerImplementation struct{}

// PostApiV1Bookings создает новую бронь
func (s *ServerImplementation) PostApiV1Bookings(ctx echo.Context) error {
	var req CreateBookingRequest

	// Парсим тело запроса
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Создаём объект для базы данных
	booking := db.BookingModel{
		CheckInDate:       req.CheckInDate.Time,
		ClientFullName:    *req.ClientFullName,
		ClientPhoneNumber: *req.ClientPhoneNumber,
		Duration:          *req.Duration,
		HotelName:         *req.HotelName,
		RoomNumber:        *req.RoomNumber,
		TotalPrice:        0.0, // Вы можете добавить логику для расчёта цены
	}

	// Сохраняем запись в базу данных
	if err := db.DB.Create(&booking).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create booking"})
	}

	// Возвращаем созданное бронирование
	return ctx.JSON(http.StatusCreated, Booking{
		BookingId:         &booking.ID,
		CheckInDate:       req.CheckInDate,
		ClientFullName:    req.ClientFullName,
		ClientPhoneNumber: req.ClientPhoneNumber,
		Duration:          req.Duration,
		HotelName:         req.HotelName,
		RoomNumber:        req.RoomNumber,
		TotalPrice:        &booking.TotalPrice,
	})
}

// GetApiV1BookingsClient возвращает список бронирований для клиента
func (s *ServerImplementation) GetApiV1BookingsClient(ctx echo.Context, params GetApiV1BookingsClientParams) error {
	var bookings []db.BookingModel

	// Извлекаем записи из базы данных
	if err := db.DB.Where("client_phone_number = ?", params.PhoneNumber).Find(&bookings).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch bookings"})
	}

	// Преобразуем записи в ответ
	var response []Booking
	for _, b := range bookings {
		response = append(response, Booking{
			BookingId:         &b.ID,
			CheckInDate:       &openapitypes.Date{Time: b.CheckInDate},
			ClientFullName:    &b.ClientFullName,
			ClientPhoneNumber: &b.ClientPhoneNumber,
			Duration:          &b.Duration,
			HotelName:         &b.HotelName,
			RoomNumber:        &b.RoomNumber,
			TotalPrice:        &b.TotalPrice,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetApiV1BookingsHotel возвращает список бронирований для отеля
func (s *ServerImplementation) GetApiV1BookingsHotel(ctx echo.Context, params GetApiV1BookingsHotelParams) error {
	var bookings []db.BookingModel

	// Извлекаем записи из базы данных по ID отеля
	if err := db.DB.Where("hotel_name = ?", params.HotelId).Find(&bookings).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch bookings for the hotel"})
	}

	// Преобразуем записи в ответ
	var response []Booking
	for _, b := range bookings {
		response = append(response, Booking{
			BookingId:         &b.ID,
			CheckInDate:       &openapitypes.Date{Time: b.CheckInDate},
			ClientFullName:    &b.ClientFullName,
			ClientPhoneNumber: &b.ClientPhoneNumber,
			Duration:          &b.Duration,
			HotelName:         &b.HotelName,
			RoomNumber:        &b.RoomNumber,
			TotalPrice:        &b.TotalPrice,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}
