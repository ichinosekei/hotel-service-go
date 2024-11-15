package hotelier

import (
	"database/sql"
)

type Hotel struct {
	ID       int
	Name     string
	Location string
}

type Room struct {
	ID         int
	HotelID    int
	RoomNumber string
	Price      float64
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateHotel(name, location string) (int, error) {
	var id int
	err := s.db.QueryRow("INSERT INTO Hotels (name, location) VALUES ($1, $2) RETURNING id", name, location).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) CreateRoom(hotelID int, roomNumber string, price float64) (int, error) {
	var id int
	err := s.db.QueryRow("INSERT INTO Rooms (hotel_id, room_number, price) VALUES ($1, $2, $3) RETURNING id", hotelID, roomNumber, price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) ListHotels() ([]Hotel, error) {
	rows, err := s.db.Query("SELECT id, name, location FROM Hotels")
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
