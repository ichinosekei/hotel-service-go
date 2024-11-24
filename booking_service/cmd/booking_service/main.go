package main

import (
	"booking_service/internal/pkg/app"
	"booking_service/internal/pkg/config"
	"booking_service/internal/pkg/persistent/repository"
	"booking_service/internal/pkg/persistent/server"
	"log"
)

func main() {
	srv_cfg := server.Config{
		URL: ":8080",
	}

	database_cfg := repository.Config{
		DSN: "host=localhost user=booking_user password=booking_password dbname=booking_db port=5432 sslmode=disable",
	}

	cfg := config.Config{
		database_cfg,
		srv_cfg,
	}

	service := app.NewBookingService()
	err := service.Init(cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = service.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server is running on port 8080")
}
