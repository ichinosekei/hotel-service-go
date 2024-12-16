package repository

import (
	"booking_service/internal/pkg/clients"
	kafka "booking_service/internal/pkg/clients/kafka"
	"booking_service/pkg/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Repository struct {
	database    *gorm.DB
	hotelAddr   string
	paymentAddr string
}

func NewRepository() *Repository {
	return &Repository{}
}
func (repo *Repository) Init(cfg Config) error {
	var err error
	repo.hotelAddr = cfg.HotelAddr
	repo.paymentAddr = cfg.PaymentAddr
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

func (repo *Repository) Create(bookingRequest *models.BookingRequest) error {
	booking := LoadBookingRequest(bookingRequest)
	booking.BookingId = uuid.NewString()

	hotelierClient, err := clients.NewHotelClient(repo.hotelAddr)
	if err != nil {
		log.Printf("Failed to create hoteler client: %v", err)
		return err
	}
	roomPrice, err := hotelierClient.GetRoomPrice(booking.HotelId, booking.RoomNumber)
	if err != nil {
		log.Printf("Failed to get hotel room price: %v", err)
		return err
	}
	booking.TotalPrice = roomPrice * ((booking.CheckOutDate.Sub(booking.CheckInDate)).Hours()) / 24
	booking.PaymentStatus = "not paid"
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
		log.Printf("Failed to checking for overlapping data: %v", err)
		return err
	}
	if len(existingBookings) > 0 {
		err := fmt.Errorf("booking dates overlap with an existing booking for room %v in hotel %v", booking.RoomNumber, booking.HotelId)
		log.Printf("Failed to create booking in data base: %v", err)
		return err
	}

	paymentClient := clients.NewPaymentClient(repo.paymentAddr, booking.TotalPrice, booking.BookingId)
	err = paymentClient.InitiatePayment()
	if err != nil {
		log.Printf("Failed to initiate payment: %v", err)
		return err
	}
	err = repo.database.Create(booking).Error
	if err != nil {
		log.Printf("Failed to create booking in data base: %v", err)
		return err
	}
	log.Printf("Booking created in data base: %v", booking.BookingId)
	return nil
}
func (repo *Repository) GetClient(phoneNumber string) (*models.Bookings, error) {
	var bookings Bookings
	err := repo.database.Where("client_phone_number = ?", phoneNumber).Find(&bookings).Error
	if err != nil {
		log.Printf("Failed to getting from data base: %v", err)
		return nil, err
	}
	var modelsBookings models.Bookings
	for _, booking := range bookings {
		modelsBookings = append(modelsBookings, *toModelsBooking(&booking))
	}
	log.Printf("Getting bookings from data base with phoneNumber: %v", phoneNumber)
	return &modelsBookings, err
}
func (repo *Repository) GetHotel(id int) (*models.Bookings, error) {
	var bookings Bookings
	err := repo.database.Where("hotel_id = ?", id).Find(&bookings).Error
	if err != nil {
		log.Printf("Failed to getting from data base: %v", err)
		return nil, err
	}
	var modelsBookings models.Bookings
	for _, booking := range bookings {
		modelsBookings = append(modelsBookings, *toModelsBooking(&booking))
	}
	log.Printf("Getting bookings from data base with hotelId: %v", id)
	return &modelsBookings, err
}

func (repo *Repository) UpdatePaymentStatusPaid(bookingId string) error {
	var booking Booking
	err := repo.database.Where("booking_id = ?", bookingId).First(&booking).Error
	if err != nil {
		log.Printf("Failed to find in data base: %v", err)
		return err
	}
	booking.PaymentStatus = "paid"
	err = repo.database.Save(&booking).Error
	if err != nil {
		log.Printf("Failed to update booking in data base: %v", err)
		return err
	}
	producer, err := kafka.NewProducer([]string{"kafka:9092"}, "booking-events")
	if err != nil {
		log.Printf("Failed to create kafka producer: %v", err)
		return err
	}
	msg, err := kafka.NewBookingEventMessages(toModelsBooking(&booking))
	if err != nil {
		log.Printf("Failed to create kafka booking event messages: %v", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = producer.Send(ctx, *msg)
	if err != nil {
		log.Printf("Failed to send message to kafka: %v", err)
	}
	err = producer.Close()
	if err != nil {
		log.Printf("Failed to close producer: %v", err)
	}
	log.Printf("Updated booking in data base: %v", booking)
	return nil
}
