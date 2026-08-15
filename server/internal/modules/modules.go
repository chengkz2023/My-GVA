package modules

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	businessexample "github.com/chengkz2023/My-GVA/server/internal/modules/business/example"
	systemmodule "github.com/chengkz2023/My-GVA/server/internal/modules/system"
)

func HTTPModules(c *container.Container) []apphttp.Module {
	return []apphttp.Module{
		businessexample.NewModule(c),
		systemmodule.NewModule(c),
	}
}
