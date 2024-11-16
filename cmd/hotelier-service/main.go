// package hotelier_service
package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/trace"
	"hotel-service-go/internal/hotelier"
	"hotel-service-go/tracing"
)

// Переменные для сервиса и трассировки
var service *hotelier.Service
var tracer trace.Tracer

func main() {
	// Устанавливаем подключение к базе данных
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Инициализация сервиса и трассировки
	service = hotelier.NewService(db)
	tracer = tracing.StartTracing("hotelier-service")

	if err := service.InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Настраиваем маршрутизатор
	r := mux.NewRouter()

	// Маршруты для работы с отелями
	r.HandleFunc("/hotels", CreateHotelHandler).Methods("POST")
	r.HandleFunc("/hotels", GetHotelsHandler).Methods("GET")
	r.HandleFunc("/hotels/{id:[0-9]+}", UpdateHotelHandler).Methods("PUT")
	r.HandleFunc("/hotels/{id:[0-9]+}", DeleteHotelHandler).Methods("DELETE")

	// Маршруты для работы с комнатами
	r.HandleFunc("/rooms", CreateRoomHandler).Methods("POST")
	r.HandleFunc("/rooms/{id:[0-9]+}", UpdateRoomHandler).Methods("PUT")
	r.HandleFunc("/rooms/{id:[0-9]+}", DeleteRoomHandler).Methods("DELETE")

	// Запускаем сервер
	log.Println("Starting hotelier service on :8080")
	http.ListenAndServe(":8080", r)
}

// Устанавливаем подключение к базе данных
func setupDatabase() (*sql.DB, error) {
	// return sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=hotelier sslmode=disable")
	return sql.Open("postgres", "host=deployments-database-1 port=5432 user=user password=password dbname=hotelier sslmode=disable")

}

// Обработчики для работы с отелями

func CreateHotelHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "CreateHotelHandler")
	defer span.End()

	var data struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := service.CreateHotel(data.Name, data.Location)
	if err != nil {
		http.Error(w, "Failed to create hotel", http.StatusInternalServerError)
		return
	}

	log.Printf("Hotel created: %d, name: %s, location: %s", id, data.Name, data.Location)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func GetHotelsHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "GetHotelsHandler")
	defer span.End()

	location := r.URL.Query().Get("location")
	hotels, err := service.GetHotels(location)
	if err != nil {
		http.Error(w, "Failed to retrieve hotels", http.StatusInternalServerError)
		return
	}

	log.Printf("Retrieved %d hotels for location: %s", len(hotels), location)
	json.NewEncoder(w).Encode(hotels)
}

func UpdateHotelHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "UpdateHotelHandler")
	defer span.End()

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var data struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := service.UpdateHotel(id, data.Name, data.Location); err != nil {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	log.Printf("Hotel updated: %d, name: %s, location: %s", id, data.Name, data.Location)
	w.WriteHeader(http.StatusOK)
}

func DeleteHotelHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "DeleteHotelHandler")
	defer span.End()

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := service.DeleteHotel(id); err != nil {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	log.Printf("Hotel deleted: %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// Обработчики для работы с комнатами

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "CreateRoomHandler")
	defer span.End()

	var data struct {
		HotelID    int     `json:"hotel_id"`
		RoomNumber string  `json:"room_number"`
		Price      float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := service.CreateRoom(data.HotelID, data.RoomNumber, data.Price)
	if err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	log.Printf("Room created: %d, hotel_id: %d, room_number: %s, price: %.2f", id, data.HotelID, data.RoomNumber, data.Price)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func UpdateRoomHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "UpdateRoomHandler")
	defer span.End()

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var data struct {
		RoomNumber string  `json:"room_number"`
		Price      float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := service.UpdateRoom(id, data.RoomNumber, data.Price); err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	log.Printf("Room updated: %d, room_number: %s, price: %.2f", id, data.RoomNumber, data.Price)
	w.WriteHeader(http.StatusOK)
}

func DeleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "DeleteRoomHandler")
	defer span.End()

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := service.DeleteRoom(id); err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	log.Printf("Room deleted: %d", id)
	w.WriteHeader(http.StatusNoContent)
}
