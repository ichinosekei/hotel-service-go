package repository

type Booking struct {
	BookingId         string  `gorm:"primaryKey"`
	CheckInDate       string  `gorm:"not null"`
	ClientFullName    string  `gorm:"not null"`
	ClientPhoneNumber string  `gorm:"not null"`
	Duration          int     `gorm:"not null"`
	HotelId           int     `gorm:"not null"`
	RoomNumber        int     `gorm:"not null"`
	TotalPrice        float32 `gorm:"not null"`
}

type Bookings []Booking
