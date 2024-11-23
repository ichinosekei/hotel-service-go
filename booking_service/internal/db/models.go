package db

import "time"

// Booking представляет таблицу бронирований
type BookingModel struct {
	ID                int       `gorm:"primaryKey"`
	CheckInDate       time.Time `gorm:"not null"`
	ClientFullName    string    `gorm:"not null"`
	ClientPhoneNumber string    `gorm:"not null"`
	Duration          int       `gorm:"not null"`
	HotelName         string    `gorm:"not null"`
	RoomNumber        string    `gorm:"not null"`
	TotalPrice        float32   `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
