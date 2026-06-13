package application

import (
	"context"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/domain"
	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, query domain.ListQuery) (ListResponse, error) {
	if s.repo == nil {
		return ListResponse{List: []RecordResponse{}}, nil
	}
	result, err := s.repo.List(ctx, query)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return ListResponse{List: []RecordResponse{}}, nil
	}
	if err != nil {
		return ListResponse{}, apperrors.New(apperrors.Internal, 0, "list records failed", err)
	}
	items := make([]RecordResponse, 0, len(result.List))
	for _, r := range result.List {
		items = append(items, fromDomain(r))
	}
	return ListResponse{
		List: items, Total: result.Total, Page: result.Page, PageSize: result.PageSize,
	}, nil
}

func (s *Service) FindByID(ctx context.Context, id uint) (RecordResponse, error) {
	if s.repo == nil {
		return RecordResponse{}, apperrors.WithMessage(apperrors.Internal, "repository unavailable")
	}
	record, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return RecordResponse{}, apperrors.New(apperrors.Internal, 0, "find record failed", err)
	}
	return fromDomain(record), nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "repository unavailable")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) DeleteByIds(ctx context.Context, ids []int) error {
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "repository unavailable")
	}
	return s.repo.DeleteByIds(ctx, ids)
}

func fromDomain(r domain.Record) RecordResponse {
	return RecordResponse{
		ID: r.ID, CreatedAt: r.CreatedAt, IP: r.IP, Method: r.Method, Path: r.Path,
		Status: r.Status, Latency: r.Latency, Agent: r.Agent,
		ErrorMessage: r.ErrorMessage, Body: r.Body, Resp: r.Resp,
		UserID: r.UserID, Username: r.Username, NickName: r.NickName,
		User: UserInRecord{ID: uint(r.UserID), UserName: r.Username, NickName: r.NickName},
	}
}

var _ pagination.Page = pagination.Page{}

func init() { _ = pagination.Page{Page: 1, PageSize: 10} }
