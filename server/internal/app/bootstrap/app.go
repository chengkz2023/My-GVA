package bootstrap

import (
	"github.com/casbin/casbin/v2"
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	casbinauthz "github.com/chengkz2023/My-GVA/server/internal/platform/authz/casbin"
	"github.com/chengkz2023/My-GVA/server/internal/platform/config"
	"github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"github.com/chengkz2023/My-GVA/server/internal/platform/logger"
	rolemysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/role/infrastructure/mysql"
	"github.com/chengkz2023/My-GVA/server/internal/platform/timer"
	"go.uber.org/zap"
)

func Initialize() *container.Container {
	vp, cfg := config.Load()
	cache := InitBlackCache(cfg.JWT)
	log := logger.New(cfg.Zap)
	zap.ReplaceGlobals(log)
	log.Info("connecting to database...")
	db, err := database.Open(cfg.Mysql, log)
	if err != nil {
		log.Error("database open failed", zap.Error(err))
	}
	var enforcer *casbin.SyncedCachedEnforcer
	if db != nil {
		log.Info("database connected, running migrations...")
		RegisterTables(db, log, cfg.System.DisableAutoMigrate)
		log.Info("migrations complete, loading seed data...")
		EnsureSystemSeedData(db, log)
		log.Info("seed data complete, initializing casbin...")
		e, err := casbinauthz.NewEnforcer(db)
		if err != nil {
			log.Fatal("casbin init failed", zap.Error(err))
		}
		enforcer = e
		log.Info("casbin initialized")
	}

	t := timer.NewTimerTask()
	InitTimer(db, t)

	c := container.New(cfg, log, db, casbinauthz.NewAuthorizer(enforcer))
	c.BlackCache = cache
	c.VP = vp
	c.Timer = t
	c.Enforcer = enforcer

	if db != nil {
		c.AuthorityChecker = rolemysql.NewRepository(db)
	}

	return c
}
