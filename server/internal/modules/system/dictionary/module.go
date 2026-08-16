package dictionary

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/domain"
	dictionarymysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/infrastructure/mysql"
	dictionaryhttp "github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/transport/http"
)

// Module 字典管理：集中维护枚举型键值（字典分类 + 字典项），供业务模块引用。
type Module struct {
	handler *dictionaryhttp.Handler
}

func NewModule(c *container.Container) *Module {
	var repo domain.Repository
	if c != nil {
		repo = dictionarymysql.NewRepository(c.DB)
	}
	service := application.NewService(repo)
	return &Module{handler: dictionaryhttp.NewHandler(service)}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	m.handler.Register(routes.Authenticated)
}
