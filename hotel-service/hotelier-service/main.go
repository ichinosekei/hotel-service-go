// package hotelier_service
package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
	"hotel-service-go/hotel-service/internal/hotelier"
	"hotel-service-go/tracing"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/url"
	"strconv"
)

// Переменные для сервиса и трассировки
var service *hotelier.Service
var tracer trace.Tracer

func main() {
	// Загружаем конфигурацию
	config := loadConfig()
	// Устанавливаем подключение к базе данных
	db, err := setupDatabase(config)
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
	r.HandleFunc("/rooms", GetRoomsHandler).Methods("GET")
	r.HandleFunc("/rooms", CreateRoomHandler).Methods("POST")
	r.HandleFunc("/rooms/{id:[0-9]+}", UpdateRoomHandler).Methods("PUT")
	r.HandleFunc("/rooms/{id:[0-9]+}", DeleteRoomHandler).Methods("DELETE")

	// Запускаем сервер
	log.Printf("Starting hotelier service on :%s", config.Server.Port)
	http.ListenAndServe(":"+config.Server.Port, r)
}

// Загружаем конфигурацию
func loadConfig() Config {
	//viper.SetConfigName("hotelier-config") // Имя файла конфигурации без расширения
	//viper.SetConfigType("yaml")            // Формат файла
	//viper.AddConfigPath("../../configs")   // Каталог для поиска файла конфигурации
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	//if err := viper.ReadInConfig(); err != nil {
	//	dir, err := os.Getwd()
	//	if err != nil {
	//		log.Fatalf("Ошибка получения текущей директории: %s", err)
	//	}
	//	log.Printf("Текущая рабочая директория: %s", dir)
	//	log.Fatalf("Error reading config file, %s", err)
	//}
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return config
}

// Устанавливаем подключение к базе данных
func setupDatabase(config Config) (*sql.DB, error) {
	connStr := config.Database.ConnectionString()
	return sql.Open("postgres", connStr)
}

// Структура конфигурации
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// Формирует строку подключения к базе данных
func (db DatabaseConfig) ConnectionString() string {
	return "host=" + db.Host +
		" port=" + db.Port +
		" user=" + db.User +
		" password=" + db.Password +
		" dbname=" + db.DBName +
		" sslmode=" + db.SSLMode
	//" sslmode=disable"
}

// Устанавливаем подключение к базе данных
//func setupDatabase() (*sql.DB, error) {
//	// return sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=hotelier sslmode=disable")
//	return sql.Open("postgres", "host=deployments-database-1 port=5432 user=user password=password dbname=hotelier sslmode=disable")
//}

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

func GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	// Начинаем span для трассировки
	_, span := tracer.Start(r.Context(), "GetRoomsHandler")
	defer span.End()

	// Получаем параметр `hotel_id` из строки запроса
	hotelIDParam := r.URL.Query().Get("hotel_id")
	var hotelID int
	var err error
	if hotelIDParam != "" {
		hotelID, err = strconv.Atoi(hotelIDParam)
		if err != nil {
			http.Error(w, "Invalid hotel_id parameter", http.StatusBadRequest)
			return
		}
	}

	// Получаем список комнат
	rooms, err := service.GetRooms(hotelID)
	if err != nil {
		http.Error(w, "Failed to retrieve rooms", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}
