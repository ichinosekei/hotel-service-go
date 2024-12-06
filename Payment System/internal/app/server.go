package app

import (
	"context"
	"payment_system/pkg/api/v1"
)

type Server interface {
	api.ServerInterface
	Run() error
	Shutdown(ctx context.Context) error
}
