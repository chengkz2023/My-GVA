package bootstrap

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	casbinauthz "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz/casbin"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/config"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/database"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/logger"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/transaction"
	rolemysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/infrastructure/mysql"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/timer"
	"go.uber.org/zap"
)

func Initialize() *container.Container {
	vp, cfg := config.Load()
	cache := InitBlackCache(cfg.JWT)
	log := logger.New(cfg.Zap)
	zap.ReplaceGlobals(log)
	log.Info("connecting to database...")
	db := database.Open(cfg.Mysql, log)
	if db != nil {
		log.Info("database connected, running migrations...")
		RegisterTables(db, log, cfg.System.DisableAutoMigrate)
		log.Info("migrations complete, loading seed data...")
		EnsureSystemSeedData(db, log)
		log.Info("seed data complete, initializing casbin...")
		_ = casbinauthz.GetEnforcer(db)
		log.Info("casbin initialized")
	}

	t := timer.NewTimerTask()
	InitTimer(db, t)

	c := container.New(cfg, log, db, transaction.NewGormManager(db), casbinauthz.NewAuthorizer())
	c.BlackCache = cache
	c.VP = vp
	c.Timer = t

	if db != nil {
		c.AuthorityChecker = rolemysql.NewRepository(db)
	}

	SetupHandlers(c)
	return c
}
