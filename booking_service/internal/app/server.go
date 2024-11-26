package app

import (
	"booking_service/pkg/api/v1"
	"context"
)

type Server interface {
	api.ServerInterface
	Run() error
	Shutdown(ctx context.Context) error
}
