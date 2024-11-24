package app

import "booking_service/internal/app/server_gen"

type Server interface {
	server_gen.ServerInterface
	Run() error
}
