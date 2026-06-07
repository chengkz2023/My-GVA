package bootstrap

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	casbinauthz "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz/casbin"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/config"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/database"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/logger"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/transaction"
	rolemysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/infrastructure/mysql"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
)

func Initialize() *container.Container {
	global.GVA_VP, _ = config.Load()
	OtherInit()
	global.GVA_LOG = logger.New()
	zap.ReplaceGlobals(global.GVA_LOG)
	global.GVA_DB = database.Open()
	Timer()
	SetupHandlers()
	if global.GVA_DB != nil {
		RegisterTables()
		EnsureSystemSeedData()
		_ = utils.GetCasbin()
	}

	c := container.New(global.GVA_CONFIG, global.GVA_LOG, global.GVA_DB, transaction.NewGormManager(global.GVA_DB), casbinauthz.NewAuthorizer())
	if global.GVA_DB != nil {
		c.AuthorityChecker = rolemysql.NewRepository(global.GVA_DB)
	}
	return c
}
