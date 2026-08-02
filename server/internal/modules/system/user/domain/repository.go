package domain

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrRepositoryUnavailable = errors.New("user repository unavailable")
)

type Repository interface {
	FindByID(ctx context.Context, id uint) (User, error)
	List(ctx context.Context, query ListQuery) (pagination.Result[User], error)
	FindPasswordHashByID(ctx context.Context, id uint) (string, error)
	UpdatePasswordHash(ctx context.Context, id uint, passwordHash string) error
	UpdateProfile(ctx context.Context, id uint, profile ProfilePatch) (User, error)
	Create(ctx context.Context, input CreateUserInput) (User, error)
	Delete(ctx context.Context, id uint) error
	SetAuthorities(ctx context.Context, input SetAuthoritiesInput) error
}

type ListQuery struct {
	Page     pagination.Page
	Username string
	NickName string
	Phone    string
	Email    string
}

type ProfilePatch struct {
	NickName  string
	HeaderImg string
	Phone     string
	Email     string
}

type CreateUserInput struct {
	Username     string
	PasswordHash string
	NickName     string
	HeaderImg    string
	AuthorityID  uint
	AuthorityIDs []uint
	Enable       int
	Phone        string
	Email        string
}

type SetAuthoritiesInput struct {
	UserID       uint
	AuthorityIDs []uint
}
