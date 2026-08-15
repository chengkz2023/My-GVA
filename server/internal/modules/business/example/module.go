package example

import (
	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/modules/business/example/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/business/example/infrastructure/memory"
	examplehttp "github.com/chengkz2023/My-GVA/server/internal/modules/business/example/transport/http"
)

// Module 是开发示例模块（不是真实业务）：
//   - 完整 DDD 分层：domain / application / infrastructure / transport
//   - 使用内存仓储，无需数据库即可运行
//   - 复制本目录即可作为新模块模板，详见 server/docs/how-to-add-module.md
type Module struct {
	handler *examplehttp.Handler
}

func NewModule(c *container.Container) *Module {
	_ = c // 示例用内存仓储，不需要容器里的 DB；真实业务请用 c.DB 构造 mysql 仓储
	repo := memory.NewRepository()
	service := application.NewService(repo)
	return &Module{handler: examplehttp.NewHandler(service)}
}

func (m *Module) RegisterHTTP(routes apphttp.Routes) {
	// 示例挂在 Public，方便免登录调试；真实业务模块请挂 routes.Authenticated
	m.handler.Register(routes.Public)
}
