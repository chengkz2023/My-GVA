package domain

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

var (
	ErrRepositoryUnavailable = errors.New("operation record repository unavailable")
	ErrRecordNotFound        = errors.New("operation record not found")
)

type ListQuery struct {
	Page   pagination.Page
	Method string
	Path   string
	Status int
}

type Repository interface {
	List(ctx context.Context, query ListQuery) (pagination.Result[Record], error)
	FindByID(ctx context.Context, id uint) (Record, error)
	Delete(ctx context.Context, id uint) error
	DeleteByIds(ctx context.Context, ids []int) error
}
