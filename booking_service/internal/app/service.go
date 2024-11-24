package app

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
