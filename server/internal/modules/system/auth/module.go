package auth

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
)

type Module struct {
	handler *Handler
}

func NewModule(c *container.Container) *Module {
	jwt := platformauth.NewJWT(platformauth.JWTConfig{
		SigningKey:  c.Config.JWT.SigningKey,
		ExpiresTime: c.Config.JWT.ExpiresTime,
		BufferTime:  c.Config.JWT.BufferTime,
		Issuer:      c.Config.JWT.Issuer,
	})
	hasher := platformauth.NewBcryptPasswordHasher()
	service := NewService(c.DB, jwt, hasher)
	return &Module{
		handler: NewHandler(service, c.Config.Captcha, c.Logger),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated, routes.Public)
}
