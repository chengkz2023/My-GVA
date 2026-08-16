package domain

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

var (
	ErrDictionaryNotFound       = errors.New("dictionary not found")
	ErrDictionaryDetailNotFound = errors.New("dictionary detail not found")
	ErrTypeExists               = errors.New("dictionary type already exists")
	ErrRepositoryUnavailable    = errors.New("dictionary repository unavailable")
)

// Dictionary 字典分类：一组同构枚举值的容器（如「性别」）。
type Dictionary struct {
	ID     uint
	Type   string
	Name   string
	Sort   int
	Status int
}

// DictionaryDetail 字典项：字典中的一个枚举值（如「男」）。
type DictionaryDetail struct {
	ID           uint
	DictionaryID uint
	Label        string
	Value        string
	Sort         int
	Status       int
}

// DictionaryWithDetails 业务引用用的完整字典（分类 + 全部启用项）。
type DictionaryWithDetails struct {
	Dictionary
	Details []DictionaryDetail
}

type SaveDictionaryInput struct {
	ID     uint
	Type   string
	Name   string
	Sort   int
	Status int
}

type SaveDetailInput struct {
	ID           uint
	DictionaryID uint
	Label        string
	Value        string
	Sort         int
	Status       int
}

type ListDictionaryQuery struct {
	Page pagination.Page
	Type string
}

// Repository 定义持久化能力，实现在 infrastructure/mysql。
type Repository interface {
	List(ctx context.Context, query ListDictionaryQuery) (pagination.Result[Dictionary], error)
	FindByID(ctx context.Context, id uint) (Dictionary, error)
	Save(ctx context.Context, input SaveDictionaryInput) (uint, error)
	Delete(ctx context.Context, id uint) error
	ListDetails(ctx context.Context, dictionaryID uint) ([]DictionaryDetail, error)
	SaveDetail(ctx context.Context, input SaveDetailInput) (uint, error)
	DeleteDetail(ctx context.Context, id uint) error
	Types(ctx context.Context) ([]DictionaryWithDetails, error)
}
