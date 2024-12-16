package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"os/signal"
	"payment_system/internal/pkg/api"
	"payment_system/internal/pkg/config"
	"syscall"
	"time"
)

func main() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Fatalf("Error loading .env.dev file: %v", err)
	}

	Config := config.Config{
		Addr: "http://" + os.Getenv("BOOKING_SERVICE_HOST") + ":" + os.Getenv("BOOKING_SERVICE_PORT") + "/api/v1/webhook/payment",
	}

	server := api.NewPaymentServer()
	server.Init(Config, echo.New(), ":"+os.Getenv("SERVER_INTERNAL_PORT"))
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
