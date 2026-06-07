package transaction

import (
	"context"

	"gorm.io/gorm"
)

type contextKey struct{}

type Manager interface {
	Within(ctx context.Context, fn func(ctx context.Context) error) error
}

type GormManager struct {
	db *gorm.DB
}

func NewGormManager(db *gorm.DB) *GormManager {
	return &GormManager{db: db}
}

func (m *GormManager) Within(ctx context.Context, fn func(ctx context.Context) error) error {
	if m == nil || m.db == nil {
		return fn(ctx)
	}
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, contextKey{}, tx))
	})
}

func DB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(contextKey{}).(*gorm.DB); ok {
		return tx
	}
	return fallback
}
