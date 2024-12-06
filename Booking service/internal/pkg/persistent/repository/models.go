package repository

import "time"

type Booking struct {
	BookingId         string    `gorm:"primaryKey"`
	CheckInDate       time.Time `gorm:"not null"`
	ClientFullName    string    `gorm:"not null"`
	ClientPhoneNumber string    `gorm:"not null"`
	CheckOutDate      time.Time `gorm:"not null"`
	HotelId           int       `gorm:"not null"`
	RoomNumber        int       `gorm:"not null"`
	TotalPrice        float64   `gorm:"not null"`
	PaymentStatus     string    `gorm:"not null"`
}

type Bookings []Booking
