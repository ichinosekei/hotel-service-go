package main

import (
	"booking_service/internal/pkg/api"
	"booking_service/internal/pkg/config"
	"booking_service/internal/pkg/persistent/repository"
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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

	DataBaseConfig := repository.Config{
		DSN: "host=" + os.Getenv("DB_HOST") +
			" user=" + os.Getenv("DB_USER") +
			" password=" + os.Getenv("DB_PASSWORD") +
			" dbname=" + os.Getenv("DB_NAME") +
			" port=" + os.Getenv("DB_PORT") +
			" sslmode=disable",
	}

	Config := config.Config{
		Database: DataBaseConfig,
	}

	server := api.NewBookingServer()
	err := server.Init(Config, echo.New(), ":"+os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Service run failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shut down: %v", err)
	}
	log.Println("Server stopped")
}
