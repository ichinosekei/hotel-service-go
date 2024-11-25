package main

import (
	"booking_service/internal/pkg/app"
	"booking_service/internal/pkg/config"
	"booking_service/internal/pkg/persistent/repository"
	"booking_service/internal/pkg/persistent/server"
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	ServerConfig := server.Config{
		PORT: ":" + os.Getenv("SERVER_PORT"),
	}

	DataBaseConfig := repository.Config{
		DSN: "host=" + os.Getenv("DB_HOST") +
			" user=" + os.Getenv("DB_USER") +
			" password=" + os.Getenv("DB_PASSWORD") +
			" dbname=" + os.Getenv("DB_NAME") +
			" port=" + os.Getenv("DB_PORT") +
			" sslmode=disable",
	}

	ServiceConfig := config.Config{
		Database: DataBaseConfig,
		Server:   ServerConfig,
	}

	service := app.NewBookingService()
	err := service.Init(ServiceConfig)
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	go func() {
		if err := service.Run(); err != nil {
			log.Fatalf("Service run failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shut down: %v", err)
	}
	log.Println("Server stopped")
}
