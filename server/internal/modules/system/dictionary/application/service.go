package application

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, query ListDictionaryQuery) (ListDictionariesResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return ListDictionariesResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return emptyList(query.Page), nil
	}
	result, err := s.repo.List(ctx, domain.ListDictionaryQuery{Page: query.Page, Type: query.Type})
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return emptyList(query.Page), nil
	}
	if err != nil {
		return ListDictionariesResponse{}, apperrors.New(apperrors.Internal, 0, "list dictionaries failed", err)
	}
	return ListDictionariesResponse{
		List:     mapDictionaries(result.List),
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

func (s *Service) Create(ctx context.Context, cmd SaveDictionaryCommand) (DictionaryResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.Type == "" || cmd.Name == "" {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Validation, "字典类型与名称必填")
	}
	if s.repo == nil {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	if cmd.Status == 0 {
		cmd.Status = 1
	}
	id, err := s.repo.Save(ctx, domain.SaveDictionaryInput{
		Type: cmd.Type, Name: cmd.Name, Sort: cmd.Sort, Status: cmd.Status,
	})
	if errors.Is(err, domain.ErrTypeExists) {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Conflict, "字典类型已存在")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	if err != nil {
		return DictionaryResponse{}, apperrors.New(apperrors.Internal, 0, "create dictionary failed", err)
	}
	return DictionaryResponse{ID: id, Type: cmd.Type, Name: cmd.Name, Sort: cmd.Sort, Status: cmd.Status}, nil
}

func (s *Service) Update(ctx context.Context, cmd SaveDictionaryCommand) (DictionaryResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.ID == 0 || cmd.Type == "" || cmd.Name == "" {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Validation, "字典类型与名称必填")
	}
	if s.repo == nil {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	id, err := s.repo.Save(ctx, domain.SaveDictionaryInput{
		ID: cmd.ID, Type: cmd.Type, Name: cmd.Name, Sort: cmd.Sort, Status: cmd.Status,
	})
	if errors.Is(err, domain.ErrDictionaryNotFound) {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.NotFound, "字典不存在")
	}
	if errors.Is(err, domain.ErrTypeExists) {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Conflict, "字典类型已存在")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return DictionaryResponse{}, apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	if err != nil {
		return DictionaryResponse{}, apperrors.New(apperrors.Internal, 0, "update dictionary failed", err)
	}
	return DictionaryResponse{ID: id, Type: cmd.Type, Name: cmd.Name, Sort: cmd.Sort, Status: cmd.Status}, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, domain.ErrDictionaryNotFound) {
		return apperrors.WithMessage(apperrors.NotFound, "字典不存在")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	if err != nil {
		return apperrors.New(apperrors.Internal, 0, "delete dictionary failed", err)
	}
	return nil
}

func (s *Service) ListDetails(ctx context.Context, dictionaryID uint) ([]DictionaryDetailResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return nil, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return []DictionaryDetailResponse{}, nil
	}
	items, err := s.repo.ListDetails(ctx, dictionaryID)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return []DictionaryDetailResponse{}, nil
	}
	if err != nil {
		return nil, apperrors.New(apperrors.Internal, 0, "list dictionary details failed", err)
	}
	return mapDetails(items), nil
}

func (s *Service) CreateDetail(ctx context.Context, cmd SaveDetailCommand) (DictionaryDetailResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.DictionaryID == 0 || cmd.Label == "" || cmd.Value == "" {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.Validation, "字典项显示值与存储值必填")
	}
	if s.repo == nil {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	if cmd.Status == 0 {
		cmd.Status = 1
	}
	id, err := s.repo.SaveDetail(ctx, domain.SaveDetailInput{
		DictionaryID: cmd.DictionaryID, Label: cmd.Label, Value: cmd.Value, Sort: cmd.Sort, Status: cmd.Status,
	})
	if errors.Is(err, domain.ErrDictionaryNotFound) {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.NotFound, "字典不存在")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	if err != nil {
		return DictionaryDetailResponse{}, apperrors.New(apperrors.Internal, 0, "create dictionary detail failed", err)
	}
	return DictionaryDetailResponse{ID: id, DictionaryID: cmd.DictionaryID, Label: cmd.Label, Value: cmd.Value, Sort: cmd.Sort, Status: cmd.Status}, nil
}

func (s *Service) UpdateDetail(ctx context.Context, cmd SaveDetailCommand) (DictionaryDetailResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.ID == 0 || cmd.Label == "" || cmd.Value == "" {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.Validation, "字典项显示值与存储值必填")
	}
	if s.repo == nil {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	id, err := s.repo.SaveDetail(ctx, domain.SaveDetailInput{
		ID: cmd.ID, Label: cmd.Label, Value: cmd.Value, Sort: cmd.Sort, Status: cmd.Status,
	})
	if errors.Is(err, domain.ErrDictionaryDetailNotFound) {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.NotFound, "字典项不存在")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return DictionaryDetailResponse{}, apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	if err != nil {
		return DictionaryDetailResponse{}, apperrors.New(apperrors.Internal, 0, "update dictionary detail failed", err)
	}
	return DictionaryDetailResponse{ID: id, DictionaryID: cmd.DictionaryID, Label: cmd.Label, Value: cmd.Value, Sort: cmd.Sort, Status: cmd.Status}, nil
}

func (s *Service) DeleteDetail(ctx context.Context, id uint) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	err := s.repo.DeleteDetail(ctx, id)
	if errors.Is(err, domain.ErrDictionaryDetailNotFound) {
		return apperrors.WithMessage(apperrors.NotFound, "字典项不存在")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return apperrors.WithMessage(apperrors.Internal, "dictionary repository unavailable")
	}
	if err != nil {
		return apperrors.New(apperrors.Internal, 0, "delete dictionary detail failed", err)
	}
	return nil
}

// Types 业务引用接口：返回全部启用字典及其启用项。
func (s *Service) Types(ctx context.Context) ([]TypeResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return nil, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return []TypeResponse{}, nil
	}
	items, err := s.repo.Types(ctx)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return []TypeResponse{}, nil
	}
	if err != nil {
		return nil, apperrors.New(apperrors.Internal, 0, "get dictionary types failed", err)
	}
	result := make([]TypeResponse, 0, len(items))
	for _, item := range items {
		result = append(result, TypeResponse{
			Type:    item.Type,
			Name:    item.Name,
			Details: mapDetails(item.Details),
		})
	}
	return result, nil
}

func emptyList(page pagination.Page) ListDictionariesResponse {
	return ListDictionariesResponse{
		List:     []DictionaryResponse{},
		Total:    0,
		Page:     page.Page,
		PageSize: page.PageSize,
	}
}

func mapDictionaries(items []domain.Dictionary) []DictionaryResponse {
	result := make([]DictionaryResponse, 0, len(items))
	for _, d := range items {
		result = append(result, DictionaryResponse{
			ID: d.ID, Type: d.Type, Name: d.Name, Sort: d.Sort, Status: d.Status,
		})
	}
	return result
}

func mapDetails(items []domain.DictionaryDetail) []DictionaryDetailResponse {
	result := make([]DictionaryDetailResponse, 0, len(items))
	for _, d := range items {
		result = append(result, DictionaryDetailResponse{
			ID: d.ID, DictionaryID: d.DictionaryID, Label: d.Label, Value: d.Value, Sort: d.Sort, Status: d.Status,
		})
	}
	return result
}
