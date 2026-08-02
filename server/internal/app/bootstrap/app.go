package bootstrap

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	casbinauthz "github.com/chengkz2023/My-GVA/server/internal/platform/authz/casbin"
	"github.com/chengkz2023/My-GVA/server/internal/platform/config"
	"github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"github.com/chengkz2023/My-GVA/server/internal/platform/logger"
	"github.com/chengkz2023/My-GVA/server/internal/platform/transaction"
	rolemysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/role/infrastructure/mysql"
	"github.com/chengkz2023/My-GVA/server/utils/timer"
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
