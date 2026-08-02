package config

import (
	"context"

	"github.com/chengkz2023/My-GVA/server/config"
	platformconfig "github.com/chengkz2023/My-GVA/server/internal/platform/config"
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
