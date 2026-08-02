package role

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/role/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/role/domain"
	rolemysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/role/infrastructure/mysql"
	rolehttp "github.com/chengkz2023/My-GVA/server/internal/modules/system/role/transport/http"
	"github.com/chengkz2023/My-GVA/server/internal/platform/authz"
)

type Module struct {
	handler *rolehttp.Handler
}

func NewModule(c *container.Container) *Module {
	var (
		repo    domain.Repository
		pp      authz.PolicyProvider
		ps      authz.PolicySyncer
		strict  bool
	)
	if c != nil {
		repo = rolemysql.NewRepository(c.DB)
		strict = c.Config.System.UseStrictAuth
		if c.Authorizer != nil {
			pp, _ = c.Authorizer.(authz.PolicyProvider)
			ps, _ = c.Authorizer.(authz.PolicySyncer)
		}
	}
	service := application.NewService(repo, pp, ps, strict)
	return &Module{
		handler: rolehttp.NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated)
}
