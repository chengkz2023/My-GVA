package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	systemapi "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api"
	systemauth "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/auth"
	systemconfig "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/config"
	systemmenu "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/menu"
	operationrecord "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record"
	systemrole "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role"
	systemstatus "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/status"
	systemuser "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/user"
	systemversion "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/version"
)

type Module struct {
	children []v2http.Module
}

func NewModule(c *container.Container) *Module {
	return &Module{
		children: []v2http.Module{
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

func (m *Module) RegisterHTTP(routes v2http.Routes) {
	for _, child := range m.children {
		child.RegisterHTTP(routes)
	}
}
