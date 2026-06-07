package bootstrap

import (
	"os"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	apimysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/infrastructure/mysql"
	operationmysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/infrastructure/mysql"
	filemysql "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/infrastructure/mysql"
	"go.uber.org/zap"
)

func RegisterTables() {
	if global.GVA_CONFIG.System.DisableAutoMigrate {
		global.GVA_LOG.Info("auto-migrate is disabled, skipping table registration")
		return
	}

	db := global.GVA_DB
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
			global.GVA_LOG.Error("register table failed", zap.Error(err))
			os.Exit(0)
		}
	}

	_ = bizModel()
	global.GVA_LOG.Info("register table success")
}
