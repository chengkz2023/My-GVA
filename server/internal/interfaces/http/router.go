package http

import (
	v2middleware "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

const V2Prefix = "/v2"

type Module interface {
	RegisterHTTP(routes Routes)
}

type Routes struct {
	Public        *gin.RouterGroup
	Authenticated *gin.RouterGroup
}

type Config struct {
	JWT v2middleware.JWTConfig
}

func RegisterV2(engine *gin.Engine, cfg Config, modules ...Module) {
	v2 := engine.Group(V2Prefix)
	authenticated := engine.Group(V2Prefix)
	authenticated.Use(v2middleware.JWTAuthWithConfig(cfg.JWT)).Use(v2middleware.CasbinHandler())

	routes := Routes{
		Public:        v2,
		Authenticated: authenticated,
	}

	v2.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})
	for _, module := range modules {
		module.RegisterHTTP(routes)
	}
}
