package container

import (
	"github.com/casbin/casbin/v2"
	"github.com/chengkz2023/My-GVA/server/config"
	"github.com/chengkz2023/My-GVA/server/internal/platform/authz"
	"github.com/chengkz2023/My-GVA/server/internal/platform/timer"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	Config           config.Server
	Logger           *zap.Logger
	DB               *gorm.DB
	Authorizer       authz.Authorizer
	AuthorityChecker authz.AuthorityChecker
	Enforcer         *casbin.SyncedCachedEnforcer
	Redis            redis.UniversalClient
	BlackCache       local_cache.Cache
	VP               *viper.Viper
	Timer            timer.Timer
	Routes           gin.RoutesInfo
}

func New(config config.Server, logger *zap.Logger, db *gorm.DB, authorizer authz.Authorizer) *Container {
	return &Container{
		Config:     config,
		Logger:     logger,
		DB:         db,
		Authorizer: authorizer,
	}
}
