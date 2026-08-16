# BoyKing Admin

基于 Gin + Vue3 + Element-Plus 的模块化后台管理脚手架。

## 项目结构

```
server/                          # Go 1.24+ 后端
├── cmd/
│   ├── admin-api/main.go        # Web 服务入口
│   └── modulegen/main.go        # 模块代码生成器
├── config/                      # 配置结构体定义
├── internal/
│   ├── app/
│   │   ├── bootstrap/           # 启动编排 (app, db, http, router, seed, timer)
│   │   ├── container/           # DI 容器
│   │   └── scaffold/            # 模块生成器引擎
│   ├── interfaces/http/         # HTTP 路由 + 中间件 (jwt, casbin_rbac, error, operation)
│   ├── modules/
│   │   ├── system/              # 系统模块 (api, auth, config, menu, operation-record, role, status, user, version)
│   │   └── business/            # 业务模块 (example — 开发示例，非真实业务)
│   └── platform/                # 平台共享层 (auth, authz/casbin, buildinfo, config, database, errors, logger, pagination, response, timer)
└── migrations/mysql/            # SQL 迁移文件

web/                             # Vue3 + Element-Plus 前端
├── src/
│   ├── api/                     # API 接口定义
│   ├── pinia/modules/           # 状态管理 (user, router, app)
│   ├── router/                  # 路由配置
│   ├── view/layout/             # 布局组件 (header, aside, tabs)
│   └── view/superAdmin/         # 管理页面
└── .env.development             # 开发环境变量 (VITE_BASE_API=/api)
```

## 常用命令

```bash
# 后端
cd server
go build ./cmd/admin-api          # 构建
go run ./cmd/admin-api            # 启动 (默认 :8888)
go run ./cmd/modulegen -name xxx  # 生成新模块
go test ./internal/...            # 运行所有测试 (26 个包)
go test ./internal -run TestArchitectureBoundaries -count=1  # 架构边界检查

# 前端
cd web
npm install
npm run dev                       # 启动 (默认 :8080)，API 代理到 :8888
```

## 架构要点

- **DDD-lite 模块化单体**: 显式 DI，框架无关的领域层
- **模块注册**: `internal/modules/modules.go` → `HTTPModules()`
- **路由前缀**: `/api/` (由 `internal/interfaces/http/router.go` 中 `APIPrefix` 常量控制)
- **中间件**: JWT (参数化 `JWTAuthWithConfig`) + Casbin RBAC + GinRecovery
- **配置加载**: 命令行 `-c` > 环境变量 `GVA_CONFIG` > Gin 模式匹配 (`config.debug.yaml` > `config.yaml`)，`configs/` 子目录优先
- **架构测试**: 6 条边界规则自动检查 domain/application/transport/platform 之间的依赖
- 所有中间件已参数化，零全局依赖

## Agent skills

### Issue tracker

Issue 存放在 GitHub Issues（repo: chengkz2023/My-GVA），统一通过 `gh` CLI 读写。见 `docs/agents/issue-tracker.md`。

### Domain docs

单上下文布局：根 `CONTEXT.md` 为词汇表、`docs/adr/` 为决策记录（文件由 /domain-modeling 按需惰性创建）。见 `docs/agents/domain.md`。
