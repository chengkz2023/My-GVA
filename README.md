# BoyKing Admin

基于 `gin + vue3 + element-plus` 精简的后台脚手架，已完成模块化重构。

## 快速开始

### 后端

```bash
cd server
# 需要 MySQL + Redis，配置见 config.yaml
go run ./cmd/admin-api
```

### 前端

```bash
cd web
npm install
npm run dev
```

## 后端架构

DDD-lite 模块化单体，显式依赖注入，框架无关的领域层。

```
server/internal/
├── app/bootstrap/      启动编排（初始化、路由、数据库、定时器）
├── app/container/      DI 容器（Config, Logger, DB, Tx, Authorizer）
├── interfaces/http/    HTTP 路由 + 中间件
├── modules/system/     系统模块（auth/user/role/menu/api/operation-record/config/status/version）
├── modules/business/   业务模块（file/example）
└── platform/           平台共享层（auth/authz/casbin/config/database/errors/logger/...）
```

- 模块注册：`server/internal/modules/modules.go`
- 50 个端点，`/api/` 前缀
- JWT 核心逻辑（Claims、Token 生成/解析）位于 `platform/auth`，零全局依赖
- Casbin RBAC 认证与鉴权
- 中间件全部参数化：`JWTAuthWithConfig` / `GinRecovery(log)` / `OperationRecord(db, log)` / `NewLimiter(rdb, log, ...)`
- 6 条架构边界规则自动检查：`go test ./internal -run TestArchitectureBoundaries`

## 工具

| 命令 | 用途 |
|------|------|
| `go run ./cmd/migrate up` | 执行数据库迁移 |
| `go run ./cmd/migrate status` | 查看迁移状态 |
| `go run ./cmd/migrate up --dry-run` | 预览待执行迁移 |
| `go run ./cmd/modulegen -name <name>` | 生成新模块代码 |

## 文档

- 详细 API 列表与架构说明：`docs/backend-handoff.md`
- 架构蓝图：`docs/backend-architecture-blueprint.md`
