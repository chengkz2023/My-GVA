package bootstrap

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/database"
	"go.uber.org/zap"
)

func Reload() error {
	global.GVA_LOG.Info("reloading system config...")

	if err := global.GVA_VP.ReadInConfig(); err != nil {
		global.GVA_LOG.Error("re-read config failed", zap.Error(err))
		return err
	}

	if global.GVA_DB != nil {
		db, _ := global.GVA_DB.DB()
		if err := db.Close(); err != nil {
			global.GVA_LOG.Error("close old db failed", zap.Error(err))
			return err
		}
	}

	global.GVA_DB = database.Open()
	OtherInit()

	if global.GVA_DB != nil {
		RegisterTables()
		EnsureSystemSeedData()
	}

	Timer()
	global.GVA_LOG.Info("system config reloaded")
	return nil
}
