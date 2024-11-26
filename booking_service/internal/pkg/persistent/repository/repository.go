package repository

import (
	"booking_service/pkg/api/v1"
	"booking_service/pkg/models"
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
	booking := LoadBookingRequest(bookingRequest)
	booking.BookingId = uuid.NewString()
	// TODO implement a grpc request to the hotel service
	booking.TotalPrice = 0

	err := repo.database.Create(booking).Error
	if err != nil {
		log.Fatalf("Failed to create booking in data base: %v", err)
	}
	return err
}
func (repo *Repository) GetClient(phoneNumber string) (*models.Bookings, error) {
	var bookings Bookings
	err := repo.database.Where("client_phone_number = ?", phoneNumber).First(&bookings).Error
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
	err := repo.database.Where("hotel_id = ?", id).First(&bookings).Error
	if err != nil {
		log.Printf("Error getting from data base: %v", err)
	}
	var modelsBookings models.Bookings
	for _, booking := range bookings {
		modelsBookings = append(modelsBookings, *toModelsBooking(&booking))
	}
	return &modelsBookings, err
}
