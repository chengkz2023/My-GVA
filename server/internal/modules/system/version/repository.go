package version

import (
	"context"

	"github.com/chengkz2023/My-GVA/server/internal/platform/buildinfo"
)

type Repository interface {
	Info(ctx context.Context) Info
}

type BuildInfoRepository struct{}

func NewBuildInfoRepository() *BuildInfoRepository {
	return &BuildInfoRepository{}
}

func (r *BuildInfoRepository) Info(ctx context.Context) Info {
	info := buildinfo.Current()
	return Info{
		AppName:     info.AppName,
		Version:     info.Version,
		Description: info.Description,
	}
}
