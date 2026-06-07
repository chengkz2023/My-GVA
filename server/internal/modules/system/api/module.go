package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/domain"
	apimysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/infrastructure/mysql"
	apihttp "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/transport/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz"
)

type Module struct {
	handler *apihttp.Handler
}

func NewModule(c *container.Container) *Module {
	var (
		repo    domain.Repository
		pp      authz.PolicyProvider
		ps      authz.PolicySyncer
		strict  bool
	)
	if c != nil {
		repo = apimysql.NewRepository(c.DB)
		strict = c.Config.System.UseStrictAuth
		if c.Authorizer != nil {
			pp, _ = c.Authorizer.(authz.PolicyProvider)
			ps, _ = c.Authorizer.(authz.PolicySyncer)
		}
	}
	service := application.NewService(repo, pp, ps, strict)
	return &Module{
		handler: apihttp.NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes v2http.Routes) {
	m.handler.Register(routes.Authenticated)
}
