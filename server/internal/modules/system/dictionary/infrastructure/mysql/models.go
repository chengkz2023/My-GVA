package mysql

import "github.com/chengkz2023/My-GVA/server/internal/platform/database"

type SysDictionary struct {
	database.GVA_MODEL
	Type   string `json:"type" gorm:"uniqueIndex;comment:字典类型标识"`
	Name   string `json:"name" gorm:"comment:字典名称"`
	Sort   int    `json:"sort" gorm:"default:0;comment:排序"`
	Status int    `json:"status" gorm:"default:1;comment:状态 1启用 2禁用"`
}

func (SysDictionary) TableName() string { return "sys_dictionaries" }

type SysDictionaryDetail struct {
	database.GVA_MODEL
	DictionaryID uint   `json:"dictionaryId" gorm:"index;comment:所属字典ID"`
	Label        string `json:"label" gorm:"comment:显示值"`
	Value        string `json:"value" gorm:"comment:存储值"`
	Sort         int    `json:"sort" gorm:"default:0;comment:排序"`
	Status       int    `json:"status" gorm:"default:1;comment:状态 1启用 2禁用"`
}

func (SysDictionaryDetail) TableName() string { return "sys_dictionary_details" }
