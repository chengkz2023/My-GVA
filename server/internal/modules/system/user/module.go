package user

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/user/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/user/domain"
	usermysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/user/infrastructure/mysql"
	userhttp "github.com/chengkz2023/My-GVA/server/internal/modules/system/user/transport/http"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
)

type Module struct {
	handler *userhttp.Handler
}

func NewModule(c *container.Container) *Module {
	var repo domain.Repository
	var checker application.AuthorityChecker
	if c != nil {
		repo = usermysql.NewRepository(c.DB)
		if c.AuthorityChecker != nil {
			checker = c.AuthorityChecker.CheckAuthorityAuth
		}
	}
	service := application.NewService(repo, platformauth.NewBcryptPasswordHasher(), checker)
	return &Module{
		handler: userhttp.NewHandler(service),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated)
}
