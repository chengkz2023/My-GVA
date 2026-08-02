// 假设这是初始化逻辑的一部分

package bootstrap

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	"github.com/chengkz2023/My-GVA/server/utils"
)

func SetupHandlers(c *container.Container) {
	utils.GlobalSystemEvents.RegisterReloadHandler(func() error {
		return Reload(c)
	})
}
