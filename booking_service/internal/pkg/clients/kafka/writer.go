package clients

import (
	"booking_service/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

func NewBookingEventMessages(booking *models.Booking) (*kafka.Message, error) {

	msgContent := map[string]interface{}{
		"client_phone": booking.ClientPhoneNumber,
		"hotel_phone":  booking.HotelierPhoneNumber,
		"check_in":     booking.CheckInDate,
		"check_out":    booking.CheckOutDate,
		"room":         string(rune(booking.RoomNumber)),
		"event_id":     booking.BookingId,
	}
	msgValue, err := json.Marshal(msgContent)
	msg := &kafka.Message{
		Key:   []byte(booking.BookingId),
		Value: msgValue,
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal hotelier message: %w", err)
	}

	return msg, nil
}
