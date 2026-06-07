package operationrecord

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/application"
	recordmysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/infrastructure/mysql"
	recordhttp "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/transport/http"
)

type Module struct {
	handler *recordhttp.Handler
}

func NewModule(c *container.Container) *Module {
	var repo *recordmysql.Repository
	if c != nil {
		repo = recordmysql.NewRepository(c.DB)
	}
	service := application.NewService(repo)
	return &Module{
		handler: recordhttp.NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes v2http.Routes) {
	m.handler.Register(routes.Authenticated)
}
