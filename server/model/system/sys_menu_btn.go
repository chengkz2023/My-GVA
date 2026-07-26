package system

import "github.com/flipped-aurora/gin-vue-admin/server/model/common"

type SysBaseMenuBtn struct {
	common.GVA_MODEL
	Name          string `json:"name" gorm:"comment:按钮关键key"`
	Desc          string `json:"desc" gorm:"按钮备注"`
	SysBaseMenuID uint   `json:"sysBaseMenuID" gorm:"comment:菜单ID"`
}
