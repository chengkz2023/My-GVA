package auth

import (
	"time"

	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/ratelimit"
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

	limiter := ratelimit.NewLimiter(
		c.Config.System.LimitCountIP,
		time.Duration(c.Config.System.LimitTimeIP)*time.Second,
	)

	// 登录 cookie 有效期读配置；配置非法时回退 7 天
	maxAge := 7 * 24 * 60 * 60
	if d, err := platformauth.ParseDuration(c.Config.JWT.ExpiresTime); err == nil && d > 0 {
		maxAge = int(d.Seconds())
	}

	return &Module{
		handler: NewHandler(service, c.Config.Captcha, c.Logger, limiter, maxAge),
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated, routes.Public)
}
