package http

import (
	"github.com/casbin/casbin/v2"
	mw "github.com/chengkz2023/My-GVA/server/internal/interfaces/http/middleware"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const APIPrefix = "/api"

type Module interface {
	RegisterHTTP(routes Routes)
}

type Routes struct {
	Public        *gin.RouterGroup
	Authenticated *gin.RouterGroup
}

type Config struct {
	JWT      mw.JWTConfig
	DB       *gorm.DB
	Logger   *zap.Logger
	Enforcer *casbin.SyncedCachedEnforcer
}

func RegisterRoutes(engine *gin.Engine, cfg Config, modules ...Module) {
	public := engine.Group(APIPrefix)
	authenticated := engine.Group(APIPrefix)
	authenticated.Use(mw.JWTAuthWithConfig(cfg.JWT)).Use(mw.CasbinHandlerWithPrefix(cfg.Enforcer, APIPrefix))
	if cfg.DB != nil {
		authenticated.Use(mw.OperationRecord(cfg.DB, cfg.Logger))
	}

	routes := Routes{
		Public:        public,
		Authenticated: authenticated,
	}

	public.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})
	for _, module := range modules {
		module.RegisterHTTP(routes)
	}
}
