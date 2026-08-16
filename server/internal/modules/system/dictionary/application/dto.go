package application

import (
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

type DictionaryResponse struct {
	ID     uint   `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Sort   int    `json:"sort"`
	Status int    `json:"status"`
}

type DictionaryDetailResponse struct {
	ID           uint   `json:"id"`
	DictionaryID uint   `json:"dictionaryId"`
	Label        string `json:"label"`
	Value        string `json:"value"`
	Sort         int    `json:"sort"`
	Status       int    `json:"status"`
}

type ListDictionariesResponse struct {
	List     []DictionaryResponse `json:"list"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
}

type ListDictionaryQuery struct {
	Page pagination.Page
	Type string
}

type SaveDictionaryCommand struct {
	ID     uint
	Type   string
	Name   string
	Sort   int
	Status int
}

type SaveDetailCommand struct {
	ID           uint
	DictionaryID uint
	Label        string
	Value        string
	Sort         int
	Status       int
}

// TypeResponse 业务引用接口的返回结构：type -> 字典名 + 启用项。
type TypeResponse struct {
	Type    string                    `json:"type"`
	Name    string                    `json:"name"`
	Details []DictionaryDetailResponse `json:"details"`
}
