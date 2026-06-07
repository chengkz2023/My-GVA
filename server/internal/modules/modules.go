package modules

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	businessexample "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/example"
	businessfile "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file"
	systemmodule "github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system"
)

func HTTPModules(c *container.Container) []v2http.Module {
	return []v2http.Module{
		businessexample.NewModule(c),
		businessfile.NewModule(c),
		systemmodule.NewModule(c),
	}
}
