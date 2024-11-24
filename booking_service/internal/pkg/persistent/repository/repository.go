package repository

import (
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

func (repo *Repository) Create(booking *Booking) error {
	return repo.database.Create(booking).Error
}
func (repo *Repository) GetClient(phoneNumber string) (Bookings, error) {
	var bookings Bookings
	err := repo.database.Where("client_phone_number = ?", phoneNumber).First(&bookings).Error
	return bookings, err
}
func (repo *Repository) GetHotel(id int) (Bookings, error) {
	var bookings Bookings
	err := repo.database.Where("hotel_id = ?", id).First(&bookings).Error
	return bookings, err
}
