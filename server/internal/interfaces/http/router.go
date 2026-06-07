package http

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
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

func RegisterV2(engine *gin.Engine, modules ...Module) {
	v2 := engine.Group(V2Prefix)
	authenticated := engine.Group(V2Prefix)
	authenticated.Use(v2middleware.JWTAuthWithConfig(v2middleware.JWTConfig{
		ExpiresTime: global.GVA_CONFIG.JWT.ExpiresTime,
		BufferTime:  global.GVA_CONFIG.JWT.BufferTime,
		SigningKey:  global.GVA_CONFIG.JWT.SigningKey,
		BlacklistCheck: func(token string) bool {
			_, ok := global.BlackCache.Get(token)
			return ok
		},
	})).Use(v2middleware.CasbinHandler())

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
