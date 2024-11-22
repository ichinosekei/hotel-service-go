package hotelier_tests

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"hotel-service-go/hotel-service/internal/hotelier"
	"testing"
)

func setupTestDB() *sql.DB {
	// Настройка тестовой базы данных в памяти или временной директории.
	db, err := sql.Open("sqlite3", ":memory:") // Пример с SQLite для тестов
	if err != nil {
		panic(err)
	}

	// Выполняем миграции (создание таблиц)
	_, _ = db.Exec(`CREATE TABLE Hotels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		location TEXT NOT NULL
	);`)
	_, _ = db.Exec(`CREATE TABLE Rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		hotel_id INTEGER NOT NULL,
		room_number TEXT NOT NULL,
		price REAL NOT NULL,
		FOREIGN KEY(hotel_id) REFERENCES Hotels(id) ON DELETE CASCADE
	);`)

	return db
}

func TestCreateHotel(t *testing.T) {
	service := hotelier.NewService(setupTestDB())

	hotelID, err := service.CreateHotel("Test Hotel", "City")
	if err != nil || hotelID == 0 {
		t.Fatalf("Expected to create hotel, got error: %v", err)
	}
}

func TestUpdateHotel(t *testing.T) {
	service := hotelier.NewService(setupTestDB())

	hotelID, _ := service.CreateHotel("Old Name", "Old City")
	err := service.UpdateHotel(hotelID, "New Name", "New City")
	if err != nil {
		t.Fatalf("Failed to update hotel: %v", err)
	}

	// Проверяем, что обновление прошло успешно
	row := service.Db.QueryRow("SELECT name, location FROM Hotels WHERE id = ?", hotelID)
	var name, location string
	_ = row.Scan(&name, &location)
	if name != "New Name" || location != "New City" {
		t.Fatalf("Hotel not updated properly. Got: %s, %s", name, location)
	}
}

func TestDeleteHotel(t *testing.T) {
	service := hotelier.NewService(setupTestDB())

	hotelID, _ := service.CreateHotel("Test Hotel", "City")
	err := service.DeleteHotel(hotelID)
	if err != nil {
		t.Fatalf("Failed to delete hotel: %v", err)
	}

	// Проверяем, что отель удалён
	row := service.Db.QueryRow("SELECT id FROM Hotels WHERE id = ?", hotelID)
	var id int
	if err := row.Scan(&id); err == nil {
		t.Fatalf("Hotel not deleted, found ID: %d", id)
	}
}

func TestCreateRoom(t *testing.T) {
	service := hotelier.NewService(setupTestDB())

	hotelID, _ := service.CreateHotel("Test Hotel", "City")
	roomID, err := service.CreateRoom(hotelID, "101", 100.0)
	if err != nil || roomID == 0 {
		t.Fatalf("Expected to create room, got error: %v", err)
	}
}

func TestUpdateRoom(t *testing.T) {
	service := hotelier.NewService(setupTestDB())

	hotelID, _ := service.CreateHotel("Test Hotel", "City")
	roomID, _ := service.CreateRoom(hotelID, "101", 100.0)

	err := service.UpdateRoom(roomID, "102", 120.0)
	if err != nil {
		t.Fatalf("Failed to update room: %v", err)
	}

	// Проверяем, что обновление прошло успешно
	row := service.Db.QueryRow("SELECT room_number, price FROM Rooms WHERE id = ?", roomID)
	var roomNumber string
	var price float64
	_ = row.Scan(&roomNumber, &price)
	if roomNumber != "102" || price != 120.0 {
		t.Fatalf("Room not updated properly. Got: %s, %.2f", roomNumber, price)
	}
}

func TestDeleteRoom(t *testing.T) {
	service := hotelier.NewService(setupTestDB())

	hotelID, _ := service.CreateHotel("Test Hotel", "City")
	roomID, _ := service.CreateRoom(hotelID, "101", 100.0)

	err := service.DeleteRoom(roomID)
	if err != nil {
		t.Fatalf("Failed to delete room: %v", err)
	}

	// Проверяем, что комната удалена
	row := service.Db.QueryRow("SELECT id FROM Rooms WHERE id = ?", roomID)
	var id int
	if err := row.Scan(&id); err == nil {
		t.Fatalf("Room not deleted, found ID: %d", id)
	}
}

func TestDeleteNonExistentHotel(t *testing.T) {
	service := hotelier.NewService(setupTestDB())

	err := service.DeleteHotel(999) // Несуществующий ID
	if err == nil {
		t.Fatalf("Expected error when deleting non-existent hotel, got nil")
	}
}

func TestDeleteNonExistentRoom(t *testing.T) {
	service := hotelier.NewService(setupTestDB())

	err := service.DeleteRoom(999) // Несуществующий ID
	if err == nil {
		t.Fatalf("Expected error when deleting non-existent room, got nil")
	}
}
