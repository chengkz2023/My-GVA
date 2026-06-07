package domain

import (
	"context"
	"errors"
)

var (
	ErrRepositoryUnavailable = errors.New("role repository unavailable")
	ErrRoleIDExists          = errors.New("authority id already exists")
	ErrRoleNotFound          = errors.New("role not found")
	ErrRoleHasUsers          = errors.New("role has users")
	ErrRoleHasChildren       = errors.New("role has child roles")
)

type Repository interface {
	Tree(ctx context.Context, authorityID uint, strict bool) ([]Role, error)
	Save(ctx context.Context, input SaveRoleInput) error
	Delete(ctx context.Context, authorityID uint) (SaveRoleInput, error)
	FindByID(ctx context.Context, authorityID uint) (Role, error)
	FindMenuIDs(ctx context.Context, authorityID uint) ([]uint, error)
	CopyMenusAndButtons(ctx context.Context, oldAuthorityID, newAuthorityID uint) error
	SetDataAuthority(ctx context.Context, input DataAuthorityInput) error
	GetDataAuthorities(ctx context.Context, authorityID uint) ([]uint, error)
	GetDescendantIDs(ctx context.Context, authorityID uint) ([]uint, error)
	CheckAuthorityAuth(ctx context.Context, adminID, targetID uint) error
}
