package domain

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

var (
	ErrRepositoryUnavailable = errors.New("api repository unavailable")
	ErrApiNotFound           = errors.New("api not found")
	ErrApiDuplicate          = errors.New("api path+method already exists")
)

type SaveApiInput struct {
	ID          uint
	Path        string
	Description string
	ApiGroup    string
	Method      string
}

type ListQuery struct {
	Page        pagination.Page
	Path        string
	Description string
	ApiGroup    string
	Method      string
	OrderKey    string
	Desc        bool
}

type Repository interface {
	List(ctx context.Context, query ListQuery) (pagination.Result[Api], error)
	GetAll(ctx context.Context) ([]Api, error)
	Groups(ctx context.Context) ([]string, error)
	Save(ctx context.Context, input SaveApiInput) (uint, error)
	Delete(ctx context.Context, id uint) (SaveApiInput, error)
	DeleteByIds(ctx context.Context, ids []int) ([]SaveApiInput, error)
	FindByID(ctx context.Context, id uint) (Api, error)
	GetIgnored(ctx context.Context) ([]Api, error)
	Ignore(ctx context.Context, path, method string, flag bool) error
	BatchCreate(ctx context.Context, apis []SaveApiInput) error
	BatchDeleteByPathMethod(ctx context.Context, apis []SaveApiInput) error
}
