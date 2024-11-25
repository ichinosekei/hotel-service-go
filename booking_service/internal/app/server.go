package app

import (
	"booking_service/internal/app/server_gen"
	"context"
)

type Server interface {
	server_gen.ServerInterface
	Run() error
	Shutdown(ctx context.Context) error
}
