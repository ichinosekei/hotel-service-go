package repository

import (
	"booking_service/internal/pkg/hotel_client"
	"booking_service/pkg/api/v1"
	"booking_service/pkg/models"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Repository struct {
	database *gorm.DB
}

func NewRepository() *Repository {
	return &Repository{}
}
func (repo *Repository) Init(cfg Config) error {
	var err error
	repo.database, err = gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}
	err = repo.database.AutoMigrate(&Booking{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return err
	}
	log.Println("Database connection successfully established!")
	return nil
}

func (repo *Repository) Create(bookingRequest *api.BookingRequest) error {
	booking, err := LoadBookingRequest(bookingRequest)

	if err != nil {
		log.Printf("Failed to create booking request: %v", err)
		return err
	}
	booking.BookingId = uuid.NewString()

	hotelierClient, err := hotel_client.NewHotelClient("")
	if err != nil {
		log.Printf("Failed to create hoteler client: %v", err)
		return err
	}
	roomPrice, err := hotelierClient.GetRoomPrice(booking.HotelId, booking.RoomNumber)
	if err != nil {
		log.Printf("Failed to get hotel room price: %v", err)
		return err
	}
	booking.TotalPrice = roomPrice * ((booking.CheckOutDate.Sub(booking.CheckInDate)).Hours())
	err = hotelierClient.Close()
	if err != nil {
		log.Printf("Failed to close hoteler client: %v", err)
		return err
	}

	var existingBookings []Booking
	err = repo.database.Where("room_number = ? AND hotel_id = ? AND (check_in_date < ? AND check_out_date > ?)",
		booking.RoomNumber, booking.HotelId, booking.CheckOutDate, booking.CheckInDate).
		Find(&existingBookings).Error
	if err != nil {
		log.Printf("Error checking for overlapping data: %v", err)
		return err
	}
	if len(existingBookings) > 0 {
		err := fmt.Errorf("booking dates overlap with an existing booking for room %v in hotel %v", booking.RoomNumber, booking.HotelId)
		log.Printf("Failed to create booking in data base: %v", err)
		return err
	}

	err = repo.database.Create(booking).Error
	if err != nil {
		log.Printf("Failed to create booking in data base: %v", err)
	}
	return err
}
func (repo *Repository) GetClient(phoneNumber string) (*models.Bookings, error) {
	var bookings Bookings
	err := repo.database.Where("client_phone_number = ?", phoneNumber).Find(&bookings).Error
	if err != nil {
		log.Printf("Error getting from data base: %v", err)
	}
	var modelsBookings models.Bookings
	for _, booking := range bookings {
		modelsBookings = append(modelsBookings, *toModelsBooking(&booking))
	}
	return &modelsBookings, err
}
func (repo *Repository) GetHotel(id int) (*models.Bookings, error) {
	var bookings Bookings
	err := repo.database.Where("hotel_id = ?", id).Find(&bookings).Error
	if err != nil {
		log.Printf("Error getting from data base: %v", err)
	}
	var modelsBookings models.Bookings
	for _, booking := range bookings {
		modelsBookings = append(modelsBookings, *toModelsBooking(&booking))
	}
	return &modelsBookings, err
}
