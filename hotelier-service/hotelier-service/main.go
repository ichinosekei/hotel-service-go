// package hotelier_service
package main

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/pkg/api"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/pkg/repository"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/pkg/tracing"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"log"
	"net"
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

	grpcServer, grpcListener, err := startGRPCServer(service)
	if err != nil {
		log.Printf("Error starting gRPC server: %v", err)
		cleanup()
		os.Exit(1) // Завершаем программу
	}

	// Запуск сервера с graceful shutdown
	GracefulShutdown(server, grpcServer, grpcListener, cleanup)
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

// Graceful shutdown для HTTP и gRPC серверов
func GracefulShutdown(httpServer *http.Server, grpcServer *grpc.Server, grpcListener net.Listener, cleanup func()) {
	// Канал для получения сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем HTTP сервер в отдельной горутине
	go func() {
		log.Printf("Starting HTTP hotelier-server on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP hotelier-server: %v", err)
		}
	}()

	// Запуск gRPC сервера в отдельной горутине
	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Ожидаем сигнал завершения
	<-stop
	log.Println("Shutting down hotelier-servers...")

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Завершаем HTTP сервер
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP Hotelier-server forced to shutdown: %v", err)
	}

	// Завершаем gRPC сервер
	grpcServer.GracefulStop()
	grpcListener.Close()

	// Закрываем базу данных
	cleanup()

	log.Println("Hotelier-server exited gracefully")
}
