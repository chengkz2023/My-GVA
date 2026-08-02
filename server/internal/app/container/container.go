package container

import (
	"github.com/chengkz2023/My-GVA/server/config"
	"github.com/chengkz2023/My-GVA/server/internal/platform/authz"
	"github.com/chengkz2023/My-GVA/server/internal/platform/transaction"
	"github.com/chengkz2023/My-GVA/server/utils/timer"
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
	Tx               transaction.Manager
	Authorizer       authz.Authorizer
	AuthorityChecker authz.AuthorityChecker
	Redis            redis.UniversalClient
	BlackCache       local_cache.Cache
	VP               *viper.Viper
	Timer            timer.Timer
	Routes           gin.RoutesInfo
}

func New(config config.Server, logger *zap.Logger, db *gorm.DB, tx transaction.Manager, authorizer authz.Authorizer) *Container {
	return &Container{
		Config:     config,
		Logger:     logger,
		DB:         db,
		Tx:         tx,
		Authorizer: authorizer,
	}
}
