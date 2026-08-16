package application

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/modules/business/example/domain"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
)

// Service 演示 application 层标准写法：
//   - 依赖 domain.Repository 接口，不依赖具体实现
//   - 把领域错误映射为 platform/errors 的结构化错误
type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, id uint) (GreetingResponse, error) {
	if s.repo == nil {
		return GreetingResponse{}, apperrors.WithMessage(apperrors.Internal, "greeting repository unavailable")
	}
	g, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, domain.ErrGreetingNotFound) {
		return GreetingResponse{}, apperrors.WithMessage(apperrors.NotFound, "greeting not found")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return GreetingResponse{}, apperrors.WithMessage(apperrors.Internal, "greeting repository unavailable")
	}
	if err != nil {
		return GreetingResponse{}, apperrors.New(apperrors.Internal, 0, "get greeting failed", err)
	}
	return greetingToResponse(g), nil
}

func (s *Service) List(ctx context.Context) (ListGreetingsResponse, error) {
	if s.repo == nil {
		return ListGreetingsResponse{List: []GreetingResponse{}}, nil
	}
	items, err := s.repo.List(ctx)
	if err != nil {
		return ListGreetingsResponse{}, apperrors.New(apperrors.Internal, 0, "list greetings failed", err)
	}
	resp := make([]GreetingResponse, 0, len(items))
	for _, g := range items {
		resp = append(resp, greetingToResponse(g))
	}
	return ListGreetingsResponse{List: resp}, nil
}

// ScopedList 行级数据权限示范：仅返回归属部门在 deptIDs 内的数据。
// deptIDs 应由调用方从「当前角色的数据权限映射」解析后传入（见角色模块 SetDataAuthority）；
// 本示例不做角色解析，只演示过滤模式本身。
func (s *Service) ScopedList(ctx context.Context, deptIDs []uint) (ListGreetingsResponse, error) {
	if s.repo == nil {
		return ListGreetingsResponse{List: []GreetingResponse{}}, nil
	}
	items, err := s.repo.ScopedList(ctx, deptIDs)
	if err != nil {
		return ListGreetingsResponse{}, apperrors.New(apperrors.Internal, 0, "scoped list greetings failed", err)
	}
	resp := make([]GreetingResponse, 0, len(items))
	for _, g := range items {
		resp = append(resp, greetingToResponse(g))
	}
	return ListGreetingsResponse{List: resp}, nil
}

func (s *Service) Create(ctx context.Context, cmd CreateGreetingCommand) (CreateGreetingResponse, error) {
	if cmd.Message == "" {
		return CreateGreetingResponse{}, apperrors.WithMessage(apperrors.Validation, "message is required")
	}
	if s.repo == nil {
		return CreateGreetingResponse{}, apperrors.WithMessage(apperrors.Internal, "greeting repository unavailable")
	}
	g, err := s.repo.Create(ctx, domain.CreateInput{Message: cmd.Message, Author: cmd.Author})
	if err != nil {
		return CreateGreetingResponse{}, apperrors.New(apperrors.Internal, 0, "create greeting failed", err)
	}
	return CreateGreetingResponse{Greeting: greetingToResponse(g)}, nil
}

func greetingToResponse(g domain.Greeting) GreetingResponse {
	return GreetingResponse{
		ID:        g.ID,
		Message:   g.Message,
		Author:    g.Author,
		DeptID:    g.DeptID,
		CreatedAt: g.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
