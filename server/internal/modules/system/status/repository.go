package status

import (
	"context"

	"github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"gorm.io/gorm"
)

type Repository interface {
	Info(ctx context.Context) Info
}

type RuntimeRepository struct {
	db *gorm.DB
}

func NewRuntimeRepository(db *gorm.DB) *RuntimeRepository {
	return &RuntimeRepository{db: db}
}

func (r *RuntimeRepository) Info(ctx context.Context) Info {
	dbStatus := Dependency{
		Configured: r.db != nil,
		OK:         false,
	}
	if r.db == nil {
		dbStatus.Message = "database is not configured"
		return infoFromDatabase(dbStatus)
	}
	if err := database.Ping(ctx, r.db); err != nil {
		dbStatus.Message = err.Error()
		return infoFromDatabase(dbStatus)
	}
	dbStatus.OK = true
	dbStatus.Message = "ok"
	return infoFromDatabase(dbStatus)
}

func infoFromDatabase(dbStatus Dependency) Info {
	info := Info{
		Status: "ok",
		Checks: Checks{
			Database: dbStatus,
		},
	}
	if dbStatus.Configured && !dbStatus.OK {
		info.Status = "degraded"
		info.Warnings = append(info.Warnings, "database check failed")
	}
	return info
}
