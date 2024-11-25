package app

import "context"

type Service struct {
	server Server
	repo   Repository
}

func NewService(srv Server, repo Repository) *Service {
	return &Service{
		server: srv,
		repo:   repo,
	}
}

func (s *Service) Run() error {
	return s.server.Run()
}

func (s *Service) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
