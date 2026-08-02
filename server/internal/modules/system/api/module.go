package api

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/api/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/api/domain"
	apimysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/api/infrastructure/mysql"
	apihttp "github.com/chengkz2023/My-GVA/server/internal/modules/system/api/transport/http"
	"github.com/chengkz2023/My-GVA/server/internal/platform/authz"
	"github.com/gin-gonic/gin"
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
		handler: apihttp.NewHandler(service, func() gin.RoutesInfo { return c.Routes }),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated)
}
