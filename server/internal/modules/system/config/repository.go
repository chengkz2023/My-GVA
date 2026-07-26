package config

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	platformconfig "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/config"
)

type Repository interface {
	Info(ctx context.Context) Info
}

type RuntimeRepository struct {
	cfg config.Server
}

func NewRuntimeRepository(cfg config.Server) *RuntimeRepository {
	return &RuntimeRepository{cfg: cfg}
}

func (r *RuntimeRepository) Info(ctx context.Context) Info {
	return Info{
		Config: platformconfig.SafeSnapshot(r.cfg),
	}
}
