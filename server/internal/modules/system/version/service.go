package version

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Info(ctx context.Context) (InfoResponse, error) {
	info := s.repo.Info(ctx)
	return InfoResponse{
		AppName:     info.AppName,
		Version:     info.Version,
		Description: info.Description,
	}, nil
}
