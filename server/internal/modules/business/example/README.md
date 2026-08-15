# business/example — 开发示例模块

> ⚠️ 本目录是**开发示例/模板**，不是真实业务。复制本目录重命名即可作为新模块的起点。

## 展示内容

- 完整 DDD 分层：`domain`（框架无关实体 + Repository 接口）/ `application`（Service + DTO + 测试）/ `infrastructure/memory`（内存仓储实现）/ `transport/http`（Gin handler + 测试）
- 标准错误映射：领域错误 → `platform/errors`（NotFound → 404、Validation → 400 等）
- 依赖注入组装（`module.go`）与路由注册（Public 演示，真实业务请挂 Authenticated）
- 无需数据库即可运行（内存仓储自带两条种子数据）

## 端点（挂在 `/api` 下，Public 免登录）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/example/greetings` | 列表 |
| GET | `/api/example/greetings/:id` | 单个（不存在返回 404） |
| POST | `/api/example/greetings` | 创建 `{"message":"...","author":"..."}` |

## 如何基于它创建新模块

1. 复制本目录为 `internal/modules/business/<your-module>/`
2. 全局替换 `example` → 你的模块名（包名、import 路径、路由前缀）
3. 需要持久化时：在 `infrastructure/` 下新增 `mysql/` 实现 `domain.Repository`，并在 `module.go` 中改用 `c.DB` 构造
4. 在 `internal/modules/modules.go` 的 `HTTPModules()` 中注册

详细规范见 `server/docs/how-to-add-module.md`。
