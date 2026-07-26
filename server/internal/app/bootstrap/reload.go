package bootstrap

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/database"
	"go.uber.org/zap"
)

func Reload(c *container.Container) error {
	c.Logger.Info("reloading system config...")
	if err := c.VP.ReadInConfig(); err != nil {
		c.Logger.Error("re-read config failed", zap.Error(err))
		return err
	}
	if c.DB != nil {
		db, _ := c.DB.DB()
		if err := db.Close(); err != nil {
			c.Logger.Error("close old db failed", zap.Error(err))
			return err
		}
	}
	c.DB = database.Open(c.Config.Mysql, c.Logger)
	c.BlackCache = InitBlackCache(c.Config.JWT)
	if c.DB != nil {
		RegisterTables(c.DB, c.Logger, c.Config.System.DisableAutoMigrate)
		EnsureSystemSeedData(c.DB, c.Logger)
	}
	InitTimer(c.DB, c.Timer)
	c.Logger.Info("system config reloaded")
	return nil
}
