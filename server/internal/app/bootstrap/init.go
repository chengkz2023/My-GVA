// 假设这是初始化逻辑的一部分

package bootstrap

import (
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

func SetupHandlers(c *container.Container) {
	utils.GlobalSystemEvents.RegisterReloadHandler(func() error {
		return Reload(c)
	})
}
