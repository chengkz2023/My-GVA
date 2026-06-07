package database

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/database/internal"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Open() *gorm.DB {
	global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
	return openMySQL(global.GVA_CONFIG.Mysql)
}

func OpenByConfig(m config.Mysql) *gorm.DB {
	return openMySQL(m)
}

func Current() *gorm.DB {
	return global.GVA_DB
}

func Ping(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return gorm.ErrInvalidDB
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func openMySQL(m config.Mysql) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(),
		DefaultStringSize:         191,
		SkipInitializeWithVersion: false,
	}
	general := m.GeneralDB
	db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(general))
	if err != nil {
		panic(err)
	}
	db.InstanceSet("gorm:table_options", "ENGINE="+m.Engine)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	return db
}
