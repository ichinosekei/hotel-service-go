// package hotelier_service
package main

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
	"hotel-service-go/hotelier-service/internal/pkg/api"
	"hotel-service-go/hotelier-service/internal/pkg/repository"
	"hotel-service-go/hotelier-service/internal/pkg/tracing"
	"log"
	"net/http"
	_ "net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Переменные для сервиса и трассировки
var service *repository.Service
var tracer trace.Tracer

func main() {
	// Загружаем конфигурацию
	config := loadConfig()
	// Устанавливаем подключение к базе данных
	db, err := setupDatabase(config)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Закрытие базы данных при завершении
	cleanup := func() {
		log.Println("Closing database (hotelier) connection...")
		if err := db.Close(); err != nil {
			log.Printf("Error (hotelier) closing database: %v", err)
		}
	}

	// Инициализация сервиса и трассировки
	service = repository.NewService(db)
	tracer = tracing.StartTracing("hotelier-service")

	if err := service.InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Настраиваем маршрутизатор
	r := mux.NewRouter()

	// экземпляр API, для передачи
	apiHandler := api.NewAPIHandler(service, tracer)

	// Маршруты для работы с отелями
	r.HandleFunc("/hotels", apiHandler.CreateHotelHandler).Methods("POST")
	r.HandleFunc("/hotels", apiHandler.GetHotelsHandler).Methods("GET")
	r.HandleFunc("/hotels/{id:[0-9]+}", apiHandler.UpdateHotelHandler).Methods("PUT")
	r.HandleFunc("/hotels/{id:[0-9]+}", apiHandler.DeleteHotelHandler).Methods("DELETE")

	// Маршруты для работы с комнатами
	r.HandleFunc("/rooms", apiHandler.GetRoomsHandler).Methods("GET")
	r.HandleFunc("/rooms", apiHandler.CreateRoomHandler).Methods("POST")
	r.HandleFunc("/rooms/{id:[0-9]+}", apiHandler.UpdateRoomHandler).Methods("PUT")
	r.HandleFunc("/rooms/{id:[0-9]+}", apiHandler.DeleteRoomHandler).Methods("DELETE")

	// Настраиваем сервер
	server := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: r,
	}
	// Запуск сервера с graceful shutdown
	GracefulShutdown(server, cleanup)
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

// Устанавливаем подключение к базе данных
func (db DatabaseConfig) ConnectionString() string {
	return "host=" + db.Host +
		" port=" + db.Port +
		" user=" + db.User +
		" password=" + db.Password +
		" dbname=" + db.DBName +
		" sslmode=" + db.SSLMode
}

// Функция для запуска сервера с graceful shutdown
func GracefulShutdown(server *http.Server, cleanup func()) {
	// Канал для получения сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Starting hotelier-server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start hotelier-server: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-stop
	log.Println("Shutting down hotelier-server...")

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Завершаем сервер
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Hotelier-server forced to shutdown: %v", err)
	}
	// закрываем базы данных
	cleanup()

	log.Println("Hotelier-server exited gracefully")
}
