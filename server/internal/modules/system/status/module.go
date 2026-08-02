package status

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"gorm.io/gorm"
)

type Module struct {
	handler *Handler
}

func NewModule(c *container.Container) *Module {
	var db *gorm.DB
	if c != nil {
		db = c.DB
	}
	repo := NewRuntimeRepository(db)
	service := NewService(repo)
	return &Module{
		handler: NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Public)
}
