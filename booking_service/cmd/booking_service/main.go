package main

import (
	"booking_service/internal/pkg/app"
	"booking_service/internal/pkg/config"
	"booking_service/internal/pkg/persistent/repository"
	"booking_service/internal/pkg/persistent/server"
	"github.com/joho/godotenv"
	"log"
	"os"
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
		log.Fatal(err)
	}
	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
}
