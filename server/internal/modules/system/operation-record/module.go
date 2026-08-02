package operationrecord

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/operation-record/application"
	recordmysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/operation-record/infrastructure/mysql"
	recordhttp "github.com/chengkz2023/My-GVA/server/internal/modules/system/operation-record/transport/http"
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

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated)
}
