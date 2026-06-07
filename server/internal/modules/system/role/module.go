package role

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/domain"
	rolemysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/infrastructure/mysql"
	rolehttp "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/transport/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz"
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

func (m *Module) RegisterHTTP(routes v2http.Routes) {
	m.handler.Register(routes.Authenticated)
}
