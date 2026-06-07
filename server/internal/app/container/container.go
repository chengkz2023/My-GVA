package container

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/transaction"
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
