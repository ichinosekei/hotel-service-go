package config

import (
	"booking_service/internal/pkg/persistent/repository"
	"booking_service/internal/pkg/persistent/server"
)

type Config struct {
	Database repository.Config
	Server   server.Config
}
