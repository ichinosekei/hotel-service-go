package config

import (
	"booking_service/internal/pkg/persistent/repository"
)

type Config struct {
	Database repository.Config
}
