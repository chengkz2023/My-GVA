package system

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	systemapi "github.com/chengkz2023/My-GVA/server/internal/modules/system/api"
	systemauth "github.com/chengkz2023/My-GVA/server/internal/modules/system/auth"
	systemconfig "github.com/chengkz2023/My-GVA/server/internal/modules/system/config"
	systemmenu "github.com/chengkz2023/My-GVA/server/internal/modules/system/menu"
	operationrecord "github.com/chengkz2023/My-GVA/server/internal/modules/system/operation-record"
	systemrole "github.com/chengkz2023/My-GVA/server/internal/modules/system/role"
	systemstatus "github.com/chengkz2023/My-GVA/server/internal/modules/system/status"
	systemuser "github.com/chengkz2023/My-GVA/server/internal/modules/system/user"
	systemversion "github.com/chengkz2023/My-GVA/server/internal/modules/system/version"
)

type Module struct {
	children []apphttp.Module
}

func NewModule(c *container.Container) *Module {
	return &Module{
		children: []apphttp.Module{
			systemapi.NewModule(c),
			systemauth.NewModule(c),
			systemconfig.NewModule(c),
			systemmenu.NewModule(c),
			operationrecord.NewModule(c),
			systemrole.NewModule(c),
			systemstatus.NewModule(c),
			systemuser.NewModule(c),
			systemversion.NewModule(c),
		},
	}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	for _, child := range m.children {
		child.RegisterHTTP(routes)
	}
}
