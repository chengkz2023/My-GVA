package bootstrap

import (
	adapter "github.com/casbin/gorm-adapter/v3"
	apimysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/api/infrastructure/mysql"
	dictionarymysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/infrastructure/mysql"
	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
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
		apimysql.SysIgnoreApi{},
		dictionarymysql.SysDictionary{},
		dictionarymysql.SysDictionaryDetail{},
		platformdb.SysUser{},
		platformdb.SysBaseMenu{},
		platformdb.SysAuthority{},
		platformdb.SysOperationRecord{},
		platformdb.SysBaseMenuParameter{},
		platformdb.SysBaseMenuBtn{},
		platformdb.SysAuthorityBtn{},
		adapter.CasbinRule{},
		platformdb.JwtBlacklist{},
	}
	for _, t := range tables {
		if err := db.AutoMigrate(t); err != nil {
			log.Fatal("register table failed", zap.Error(err))
		}
	}

	log.Info("register table success")
}
