package status

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
		Status: info.Status,
		Checks: ChecksResponse{
			Database: DependencyStatus{
				Configured: info.Checks.Database.Configured,
				OK:         info.Checks.Database.OK,
				Message:    info.Checks.Database.Message,
			},
		},
		Warnings: append([]string(nil), info.Warnings...),
	}, nil
}
