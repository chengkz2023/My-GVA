package file

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/modules/business/file/application"
	filemysql "github.com/chengkz2023/My-GVA/server/internal/modules/business/file/infrastructure/mysql"
	filehttp "github.com/chengkz2023/My-GVA/server/internal/modules/business/file/transport/http"
)

type Module struct {
	handler *filehttp.Handler
}

func NewModule(c *container.Container) *Module {
	storePath := "uploads/file"
	if c != nil {
		storePath = c.Config.Local.StorePath
	}
	repo := filemysql.NewRepository(c.DB)
	service := application.NewService(repo, storePath)
	return &Module{
		handler: filehttp.NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated)
}
