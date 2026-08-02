package example

import "context"

import apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Info(ctx context.Context) (InfoResponse, error) {
	info := s.repo.Info(ctx)
	return InfoResponse{
		Name:    info.Name,
		Message: info.Message,
	}, nil
}

func (s *Service) Missing(ctx context.Context) (InfoResponse, error) {
	return InfoResponse{}, apperrors.WithMessage(apperrors.NotFound, "example not found")
}
