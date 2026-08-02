package bootstrap

import (
	"os"

	adapter "github.com/casbin/gorm-adapter/v3"
	apimysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/api/infrastructure/mysql"
	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
	operationmysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/operation-record/infrastructure/mysql"
	filemysql "github.com/chengkz2023/My-GVA/server/internal/modules/business/file/infrastructure/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterTables(db *gorm.DB, log *zap.Logger, disableAutoMigrate bool) {
	if disableAutoMigrate {
		log.Info("auto-migrate is disabled, skipping table registration")
		return
	}

	tables := []interface{}{
		apimysql.SysApi{},
		platformdb.SysUser{},
		platformdb.SysBaseMenu{},
		platformdb.JwtBlacklist{},
		platformdb.SysAuthority{},
		platformdb.SysDictionary{},
		operationmysql.SysOperationRecord{},
		platformdb.SysDictionaryDetail{},
		platformdb.SysBaseMenuParameter{},
		platformdb.SysBaseMenuBtn{},
		platformdb.SysAuthorityBtn{},
		platformdb.SysParams{},
		platformdb.SysVersion{},
		platformdb.SysError{},
		adapter.CasbinRule{},

		filemysql.ExaFileUploadAndDownload{},
	}
	for _, t := range tables {
		if err := db.AutoMigrate(t); err != nil {
			log.Error("register table failed", zap.Error(err))
			os.Exit(0)
		}
	}

	_ = bizModel(db)
	log.Info("register table success")
}
