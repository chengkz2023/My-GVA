package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common"
)

type JwtBlacklist struct {
	common.GVA_MODEL
	Jwt string `gorm:"type:text;comment:jwt"`
}
