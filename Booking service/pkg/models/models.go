package models

import (
	"time"
)

type Booking struct {
	BookingId         string
	CheckInDate       time.Time
	CheckOutDate      time.Time
	ClientFullName    string
	ClientPhoneNumber string
	HotelId           int
	RoomNumber        int
	TotalPrice        float64
}

type BookingRequest struct {
	CheckInDate       time.Time
	CheckOutDate      time.Time
	ClientFullName    string
	ClientPhoneNumber string
	HotelId           int
	RoomNumber        int
}

type Bookings []Booking
