package menu

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/menu/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/menu/domain"
	menumysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/menu/infrastructure/mysql"
	menuhttp "github.com/chengkz2023/My-GVA/server/internal/modules/system/menu/transport/http"
)

type Module struct {
	handler *menuhttp.Handler
}

func NewModule(c *container.Container) *Module {
	var repo domain.Repository
	var checker application.AuthorityChecker
	if c != nil {
		repo = menumysql.NewRepository(c.DB)
		if c.AuthorityChecker != nil {
			checker = c.AuthorityChecker.CheckAuthorityAuth
		}
	}
	service := application.NewService(repo, checker)
	return &Module{
		handler: menuhttp.NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated)
}
