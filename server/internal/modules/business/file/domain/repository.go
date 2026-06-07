package domain

import (
	"context"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
)

var ErrRepositoryUnavailable = errors.New("file repository unavailable")

type Storage interface {
	Upload(key string, data []byte) (url string, err error)
	Delete(key string) error
}

type Repository interface {
	List(ctx context.Context, query ListQuery) (pagination.Result[File], error)
	FindByID(ctx context.Context, id uint) (File, error)
	Create(ctx context.Context, file File) (File, error)
	Update(ctx context.Context, id uint, name string, tag string) (File, error)
	Delete(ctx context.Context, id uint) (File, error)
}

type ListQuery struct {
	Page    pagination.Page
	ClassID int
	Name    string
}
