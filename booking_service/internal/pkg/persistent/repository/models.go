package repository

import "time"

type Booking struct {
	ID                string    `gorm:"primaryKey"`
	CheckInDate       time.Time `gorm:"not null"`
	ClientFullName    string    `gorm:"not null"`
	ClientPhoneNumber string    `gorm:"not null"`
	Duration          int       `gorm:"not null"`
	HotelId           int       `gorm:"not null"`
	RoomNumber        int       `gorm:"not null"`
	TotalPrice        float32   `gorm:"not null"`
}

type Bookings []Booking
