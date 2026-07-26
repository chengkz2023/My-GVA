package bootstrap

import (
	"os"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	apimysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/infrastructure/mysql"
	operationmysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/infrastructure/mysql"
	filemysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/infrastructure/mysql"
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
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		operationmysql.SysOperationRecord{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysParams{},
		system.SysVersion{},
		system.SysError{},
		adapter.CasbinRule{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		filemysql.ExaFileUploadAndDownload{},
		example.ExaAttachmentCategory{},
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
