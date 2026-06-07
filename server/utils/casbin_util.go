package utils

import (
	"github.com/casbin/casbin/v2"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	casbinauthz "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz/casbin"
)

func GetCasbin() *casbin.SyncedCachedEnforcer {
	return casbinauthz.GetEnforcer(global.GVA_DB)
}
