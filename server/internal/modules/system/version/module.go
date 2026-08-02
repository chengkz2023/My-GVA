package version

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
)

type Module struct {
	handler *Handler
}

func NewModule(c *container.Container) *Module {
	_ = c
	repo := NewBuildInfoRepository()
	service := NewService(repo)
	return &Module{
		handler: NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Public)
}
