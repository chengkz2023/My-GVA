package config

import (
	"context"

	platformconfig "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/config"
)

type Repository interface {
	Info(ctx context.Context) Info
}

type RuntimeRepository struct{}

func NewRuntimeRepository() *RuntimeRepository {
	return &RuntimeRepository{}
}

func (r *RuntimeRepository) Info(ctx context.Context) Info {
	return Info{
		Config: platformconfig.SafeSnapshot(),
	}
}
