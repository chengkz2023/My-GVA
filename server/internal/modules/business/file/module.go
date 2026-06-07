package file

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/application"
	filemysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/infrastructure/mysql"
	filehttp "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/transport/http"
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

func (m *Module) RegisterHTTP(routes v2http.Routes) {
	m.handler.Register(routes.Authenticated)
}
