package example

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
)

type Module struct {
	handler *Handler
}

func NewModule(c *container.Container) *Module {
	_ = c
	repo := NewMemoryRepository()
	service := NewService(repo)
	return &Module{
		handler: NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes v2http.Routes) {
	m.handler.Register(routes.Public)
}
