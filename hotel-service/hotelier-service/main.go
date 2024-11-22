// package hotelier_service
package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/yaml.v3"
	"hotel-service-go/hotel-service/internal/hotelier"
	"hotel-service-go/hotel-service/tracing"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/url"
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
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
	return config
}

// Устанавливаем подключение к базе данных
func setupDatabase(config Config) (*sql.DB, error) {
	connStr := config.Database.ConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
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

// Устанавливаем подключение к базе данных
func (db DatabaseConfig) ConnectionString() string {
	return "host=" + db.Host +
		" port=" + db.Port +
		" user=" + db.User +
		" password=" + db.Password +
		" dbname=" + db.DBName +
		" sslmode=" + db.SSLMode
}

// Устанавливаем подключение к базе данных
//func setupDatabase() (*sql.DB, error) {
//	// return sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=hotelier sslmode=disable")
//	return sql.Open("postgres", "host=deployments-database-1 port=5432 user=user password=password dbname=hotelier sslmode=disable")
//}
