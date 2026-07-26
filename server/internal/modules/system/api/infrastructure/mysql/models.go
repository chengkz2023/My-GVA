package mysql

import "github.com/flipped-aurora/gin-vue-admin/server/model/common"

type SysApi struct {
	common.GVA_MODEL
	Path        string `json:"path" gorm:"comment:api路径"`
	Description string `json:"description" gorm:"comment:api中文描述"`
	ApiGroup    string `json:"apiGroup" gorm:"comment:api组"`
	Method      string `json:"method" gorm:"default:POST;comment:方法"`
}

func (SysApi) TableName() string { return "sys_apis" }

type SysIgnoreApi struct {
	common.GVA_MODEL
	Path   string `json:"path" gorm:"comment:api路径"`
	Method string `json:"method" gorm:"default:POST;comment:方法"`
	Flag   bool   `json:"flag" gorm:"-"`
}

func (SysIgnoreApi) TableName() string { return "sys_ignore_apis" }
