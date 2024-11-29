package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
)

type Service struct {
	Db *sql.DB
	//Tracer trace.Tracer
}

func NewService(db *sql.DB) *Service {
	return &Service{
		Db: db,
		//Tracer: otel.Tracer("hotelier-service"),
	}
}

type Hotel struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// Структура для представления комнаты
type Room struct {
	ID         int     `json:"id"`
	HotelID    int     `json:"hotel_id"`
	RoomNumber string  `json:"room_number"`
	Price      float64 `json:"price"`
}

// CreateHotel creates a new hotel.
func (s *Service) CreateHotel(name, location string) (int, error) {
	var id int
	err := s.Db.QueryRow(
		"INSERT INTO Hotels (name, location) VALUES ($1, $2) RETURNING id",
		name, location,
	).Scan(&id)
	if err != nil {
		log.Printf("Error creating hotel: %v", err)
		return 0, err
	}
	log.Printf("Hotel created with ID: %d", id)
	return id, nil
}

// UpdateHotel updates an existing hotel's information.
func (s *Service) UpdateHotel(id int, name, location string) error {
	result, err := s.Db.Exec("UPDATE Hotels SET name = $1, location = $2 WHERE id = $3", name, location, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("hotel not found")
	}
	return nil
}

// DeleteHotel deletes a hotel by ID.
func (s *Service) DeleteHotel(id int) error {
	result, err := s.Db.Exec("DELETE FROM Hotels WHERE id = $1", id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("hotel not found")
	}
	return nil
}

// CreateRoom creates a new room for a hotel.
func (s *Service) CreateRoom(hotelID int, roomNumber string, price float64) (int, error) {
	var id int
	err := s.Db.QueryRow("INSERT INTO Rooms (hotel_id, room_number, price) VALUES ($1, $2, $3) RETURNING id", hotelID, roomNumber, price).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "foreign_key_violation" {
			return 0, errors.New("hotel not found")
		}
		return 0, err

	}
	return id, nil
}

// GetRooms возвращает список комнат, возможно с фильтрацией по hotelID
func (s *Service) GetRooms(hotelID int) ([]Room, error) {
	query := "SELECT id, hotel_id, room_number, price FROM Rooms"
	args := []interface{}{}

	if hotelID > 0 {
		query += " WHERE hotel_id = $1"
		args = append(args, hotelID)
	}

	rows, err := s.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.HotelID, &room.RoomNumber, &room.Price); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

// UpdateRoom updates room information.
func (s *Service) UpdateRoom(id int, roomNumber string, price float64) error {
	result, err := s.Db.Exec("UPDATE Rooms SET room_number = $1, price = $2 WHERE id = $3", roomNumber, price, id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("room not found")
	}
	return nil
}

// DeleteRoom deletes a room by ID.
func (s *Service) DeleteRoom(id int) error {
	result, err := s.Db.Exec("DELETE FROM Rooms WHERE id = $1", id)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("room not found")
	}
	return nil
}

// GetHotels возвращает список отелей, возможно с фильтрацией по местоположению
func (s *Service) GetHotels(location string) ([]Hotel, error) {
	query := "SELECT id, name, location FROM Hotels"
	args := []interface{}{}

	if location != "" {
		query += " WHERE location = $1"
		args = append(args, location)
	}

	rows, err := s.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hotels []Hotel
	for rows.Next() {
		var hotel Hotel
		if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Location); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}
	return hotels, nil
}

func (s *Service) InitializeDatabase() error {
	if s.Db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	_, err := s.Db.Exec(`
        CREATE TABLE IF NOT EXISTS Hotels (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            location VARCHAR(100) NOT NULL
        );
        CREATE TABLE IF NOT EXISTS Rooms (
            id SERIAL PRIMARY KEY,
            hotel_id INTEGER NOT NULL REFERENCES Hotels(id) ON DELETE CASCADE,
            room_number VARCHAR(10) NOT NULL,
            price NUMERIC(10, 2) NOT NULL
        );
    `)
	return err
}
