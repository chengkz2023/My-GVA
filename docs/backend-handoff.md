# 后端模块化重构 — 完成报告

## 状态：已完成

基于 `docs/backend-architecture-blueprint.md` 的后端重构已全部完成。所有模块、平台包和基础设施均已就位，旧架构（`api/`、`router/`、`service/`、`model/`、`global/` 等）已完全移除。

## 验证命令

```bash
cd server
go build . && go build ./cmd/admin-api && go build ./cmd/modulegen
go test ./...
```

`go test ./...` 覆盖全部包，其中 26 个包包含测试文件（均位于 `./internal/...` 下；`utils/` 已移除）。

## 架构

### 目录结构

```
server/
├── cmd/
│   ├── admin-api/main.go            # Web 服务入口（调用 bootstrap.Run()）
│   └── modulegen/main.go            # 模块代码生成器
├── internal/
│   ├── app/
│   │   ├── bootstrap/               # 启动编排（app, database, http, redis, router, seed, timer）
│   │   ├── container/               # DI 容器（Config, Logger, DB, Tx, Authorizer, Redis, Timer 等）
│   │   ├── scaffold/                # 模块生成器脚手架引擎
│   │   └── task/                    # 定时维护任务（clearTable）
│   ├── interfaces/http/
│   │   ├── router.go                # RegisterRoutes() — public + authenticated 路由组
│   │   └── middleware/              # jwt, casbin_rbac, error
│   ├── modules/
│   │   ├── system/                  # auth, user, role, menu, api, operation-record, config, status, version
│   │   └── business/                # example（开发示例模块，非真实业务）
│   ├── platform/                    # 共享包（不导入 modules）
│   │   ├── auth/                    # Actor, Claims, JWT, Token, Password（零全局依赖）
│   │   ├── authz/                   # Authorizer + PolicyProvider + PolicySyncer + AuthorityChecker 接口
│   │   ├── authz/casbin/            # Casbin 实现（构造注入）
│   │   ├── buildinfo/               # 构建元信息
│   │   ├── config/                  # 基于 Viper 的配置加载
│   │   ├── database/                # GORM 数据库连接
│   │   ├── errors/                  # Kind 枚举错误类型
│   │   ├── logger/                  # Zap 日志器
│   │   ├── pagination/              # 泛型分页
│   │   ├── response/                # 统一响应（code/msg/message/data）
│   │   └── timer/                   # 定时任务调度（cron）
│   └── architecture_test.go         # 6 条架构边界规则
├── config/                          # 配置结构体定义（system, jwt, redis, gorm-mysql, zap, cors, captcha）
├── configs/                         # YAML 配置文件（config.debug.yaml, config.example.yaml）
├── migrations/mysql/                # SQL 迁移文件（含 README）
├── Dockerfile
├── go.mod / go.sum
└── README.md
```

### HTTP 路由

- `GET /api/health` — 公开健康检查
- `POST /api/login` — 登录认证
- `POST /api/base/captcha` — 图形验证码
- 其余端点均在 `/api/` 路径下，由 JWT + Casbin 中间件保护
- 无需鉴权的公开端点：`health`、`login`、`captcha`、`version/info`、`example/greetings`

### 模块结构

每个模块遵循以下两种模式之一：

**完整分层**（domain/application/infrastructure/transport）：
```
module/
├── domain/              # 框架无关的类型 + Repository 接口
├── application/         # Service + DTO + 测试
├── infrastructure/mysql/# GORM 仓储实现
├── transport/http/      # Gin handler + 测试
└── module.go            # 组装（NewModule → repo → service → handler）
```
适用模块：`business/example`（内存仓储，开发示例）、`system/user`、`system/role`、`system/menu`、`system/api`、`system/operation-record`

**轻量级**（扁平结构）：
```
module/
├── handler.go / service.go / repository.go / model.go / dto.go / module.go
└── handler_test.go
```
适用模块：`system/auth`、`system/config`、`system/status`、`system/version`

## 全部端点（50 个）

### 公开端点

| 模块 | 方法 | 路径 | 鉴权 |
|------|------|------|------|
| — | GET | `/api/health` | 否 |
| auth | POST | `/api/login` | 否 |
| auth | POST | `/api/base/captcha` | 否 |
| version | GET | `/api/system/version/info` | 否 |
| example | GET | `/api/example/greetings` | 否 |
| example | GET | `/api/example/greetings/:id` | 否 |
| example | POST | `/api/example/greetings` | 否 |

### 系统模块（需鉴权）

| 模块 | 方法 | 路径 | 鉴权 |
|------|------|------|------|
| auth | GET | `/api/system/auth/me` | 是 |
| auth | POST | `/api/system/auth/logout` | 是 |
| user | GET | `/api/system/user/me` | 是 |
| user | GET | `/api/system/user/list` | 是 |
| user | POST | `/api/system/user/password` | 是 |
| user | PUT | `/api/system/user/profile` | 是 |
| user | POST | `/api/system/user` | 是 |
| user | DELETE | `/api/system/user/:id` | 是 |
| user | POST | `/api/system/user/:id/reset-password` | 是 |
| user | PUT | `/api/system/user/:id/authorities` | 是 |
| role | GET | `/api/system/role/tree` | 是 |
| role | GET | `/api/system/role/:id` | 是 |
| role | GET | `/api/system/role/:id/data-authorities` | 是 |
| role | POST | `/api/system/role` | 是 |
| role | PUT | `/api/system/role/:id` | 是 |
| role | DELETE | `/api/system/role/:id` | 是 |
| role | POST | `/api/system/role/copy` | 是 |
| role | POST | `/api/system/role/data-authority` | 是 |
| menu | GET | `/api/system/menu/tree` | 是 |
| menu | POST | `/api/system/menu/authority` | 是 |
| menu | POST | `/api/system/menu` | 是 |
| menu | GET | `/api/system/menu/:id` | 是 |
| menu | DELETE | `/api/system/menu/:id` | 是 |
| api | GET | `/api/system/api/list` | 是 |
| api | GET | `/api/system/api/all` | 是 |
| api | GET | `/api/system/api/groups` | 是 |
| api | GET | `/api/system/api/policies/:authorityId` | 是 |
| api | PUT | `/api/system/api/policies` | 是 |
| api | GET | `/api/system/api/:id` | 是 |
| api | POST | `/api/system/api` | 是 |
| api | PUT | `/api/system/api/:id` | 是 |
| api | DELETE | `/api/system/api/:id` | 是 |
| api | POST | `/api/system/api/batch-delete` | 是 |
| api | POST | `/api/system/api/fresh-casbin` | 是 |
| api | GET | `/api/system/api/sync` | 是 |
| api | POST | `/api/system/api/ignore` | 是 |
| api | POST | `/api/system/api/batch-sync` | 是 |
| operation-record | GET | `/api/system/operation-record/list` | 是 |
| operation-record | GET | `/api/system/operation-record/:id` | 是 |
| operation-record | DELETE | `/api/system/operation-record/:id` | 是 |
| operation-record | POST | `/api/system/operation-record/batch-delete` | 是 |
| config | GET | `/api/system/config/info` | 是 |
| status | GET | `/api/system/status/info` | 是 |

### 业务模块

`business/example` 为开发示例模块（非真实业务，见 `server/internal/modules/business/example/README.md`），其端点为公开端点：`/api/example/greetings`（GET/POST）与 `/api/example/greetings/:id`（GET）。

## 中间件状态

现有 4 个中间件，全部参数化、零 `global` 依赖：

- `jwt.go` — `JWTAuthWithConfig(JWTConfig{ExpiresTime, BufferTime, SigningKey, BlacklistCheck})`
- `casbin_rbac.go` — `CasbinHandlerWithPrefix(enforcer, prefix)`（enforcer 由构造注入）
- `error.go` — `GinRecovery(log, stack)`
- `operation.go` — `OperationRecord(db, log)`（写操作审计入库）

已移除的旧中间件：`cors.go`、`limit_ip.go`（对应能力已由 platform 层或启动编排接管）。

## 已删除

旧架构目录已全部移除：`api/`、`router/`、`service/`、`source/`、`core/`、`initialize/`、`model/`、`global/`、`cmd/migrate`、`internal/app/migration/`、`utils/`、`internal/platform/validator/`、`task/`、`internal/platform/transaction/`、`internal/modules/business/file/`、`config/oss_local.go`（Local 文件存储配置）

## 剩余工作（低优先级）

- `configs/` 目前仅保留 `config.debug.yaml` 与 `config.example.yaml`，生产 `config.yaml` 待补充
- `migrations/mysql/` 目前仅含 `README.md`，版本化 SQL 迁移文件待补充

## 如何添加新模块

```bash
cd server
go run ./cmd/modulegen -name <模块名>
```

然后在 `internal/modules/modules.go` 中注册。所有模块遵循依赖规则：

"platform 不得导入 modules。domain 不得导入框架。handler 不得触碰持久化。"
