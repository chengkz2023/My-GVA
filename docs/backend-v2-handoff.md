# 后端 V2 重构 — 完成报告

## 状态：已完成

基于 `docs/backend-architecture-blueprint.md` 的 V2 后端重构已全部完成。所有模块、平台包和基础设施均已就位，旧架构已完全替换。

## 验证命令

```bash
cd server
go build . && go build ./cmd/admin-api && go build ./cmd/migrate && go build ./cmd/modulegen
go test ./...
```

全部 28 个测试包通过，4 个构建均通过。

## 架构

### 目录结构

```
server/
├── main.go                          # 入口（调用 bootstrap.Run()）
├── cmd/
│   ├── admin-api/main.go            # 未来入口点
│   ├── migrate/main.go              # 数据库迁移工具（up/down/status/dry-run）
│   └── modulegen/main.go            # 模块代码生成器
├── internal/
│   ├── app/
│   │   ├── bootstrap/               # 启动编排（初始化、路由、数据库、定时器）
│   │   ├── container/               # DI 容器（Config, Logger, DB, Tx, Authorizer）
│   │   ├── migration/               # 迁移执行引擎
│   │   └── scaffold/                # 模块生成器脚手架
│   ├── interfaces/http/
│   │   ├── router.go                # RegisterV2() — public + authenticated 路由组
│   │   └── middleware/              # jwt, casbin_rbac, cors, error, limit_ip, operation
│   ├── modules/
│   │   ├── system/                  # auth, user, role, menu, api, operation-record, config, status, version
│   │   └── business/                # file, example
│   ├── platform/                    # 共享包（不导入 modules）
│   │   ├── auth/                    # Actor, Claims, JWT, Token, Password（零全局依赖）
│   │   ├── authz/                   # Authorizer + PolicyProvider + PolicySyncer + AuthorityChecker 接口
│   │   ├── authz/casbin/            # Casbin 实现
│   │   ├── config/                  # 基于 Viper 的配置加载
│   │   ├── database/                # GORM 数据库连接
│   │   ├── errors/                  # Kind 枚举错误类型
│   │   ├── logger/                  # Zap 日志器
│   │   ├── pagination/              # 泛型分页
│   │   ├── response/                # 统一响应（code/msg/message/data）
│   │   ├── transaction/             # 事务管理器
│   │   ├── validator/               # 验证器
│   │   └── buildinfo/               # 构建元信息
│   └── architecture_test.go        # 6 条架构边界规则
├── config/                          # 配置结构体定义
├── migrations/mysql/                # SQL 迁移文件
├── model/                           # 旧版 GORM 模型（25 个文件，仍在活跃使用）
├── utils/                           # JWT, claims, timer, AST 等
└── global/                          # 旧版全局变量（GVA_DB, GVA_CONFIG 等）
```

### HTTP 路由

- `GET /v2/health` — 公开健康检查
- `POST /v2/login` — 登录认证
- 其他所有端点均在 `/v2/` 路径下，使用 JWT + Casbin 中间件保护

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

**轻量级**（扁平结构）：
```
module/
├── handler.go / service.go / repository.go / model.go / dto.go / module.go
└── handler_test.go
```

## 全部 V2 端点（50 个）

### 系统模块

| 模块 | 方法 | 路径 | 鉴权 |
|------|------|------|------|
| auth | POST | `/v2/login` | 否 |
| auth | GET | `/v2/system/auth/me` | 是 |
| user | GET | `/v2/system/user/me` | 是 |
| user | GET | `/v2/system/user/list` | 是 |
| user | GET | `/v2/system/user/:id` | 是 |
| user | POST | `/v2/system/user` | 是 |
| user | PUT | `/v2/system/user/profile` | 是 |
| user | PUT | `/v2/system/user/password` | 是 |
| user | DELETE | `/v2/system/user/:id` | 是 |
| user | PUT | `/v2/system/user/reset-password` | 是 |
| user | PUT | `/v2/system/user/set-authorities` | 是 |
| role | GET | `/v2/system/role/tree` | 是 |
| role | GET | `/v2/system/role/:id` | 是 |
| role | GET | `/v2/system/role/data-authorities` | 是 |
| role | POST | `/v2/system/role` | 是 |
| role | PUT | `/v2/system/role/:id` | 是 |
| role | DELETE | `/v2/system/role/:id` | 是 |
| role | POST | `/v2/system/role/copy` | 是 |
| role | PUT | `/v2/system/role/set-data-authority` | 是 |
| menu | GET | `/v2/system/menu/tree` | 是 |
| menu | GET | `/v2/system/menu/:id` | 是 |
| menu | POST | `/v2/system/menu` | 是 |
| menu | DELETE | `/v2/system/menu/:id` | 是 |
| menu | PUT | `/v2/system/menu/assign-authority` | 是 |
| api | GET | `/v2/system/api/list` | 是 |
| api | GET | `/v2/system/api/all` | 是 |
| api | GET | `/v2/system/api/groups` | 是 |
| api | GET | `/v2/system/api/policies/:authorityId` | 是 |
| api | PUT | `/v2/system/api/policies` | 是 |
| api | GET | `/v2/system/api/:id` | 是 |
| api | POST | `/v2/system/api` | 是 |
| api | PUT | `/v2/system/api/:id` | 是 |
| api | DELETE | `/v2/system/api/:id` | 是 |
| api | POST | `/v2/system/api/batch-delete` | 是 |
| api | POST | `/v2/system/api/fresh-casbin` | 是 |
| api | GET | `/v2/system/api/sync` | 是 |
| api | POST | `/v2/system/api/ignore` | 是 |
| api | POST | `/v2/system/api/batch-sync` | 是 |
| operation-record | GET | `/v2/system/operation-record/list` | 是 |
| operation-record | GET | `/v2/system/operation-record/:id` | 是 |
| operation-record | DELETE | `/v2/system/operation-record/:id` | 是 |
| operation-record | POST | `/v2/system/operation-record/batch-delete` | 是 |
| config | GET | `/v2/system/config/info` | 否 |
| status | GET | `/v2/system/status/info` | 否 |
| version | GET | `/v2/system/version/info` | 否 |

### 业务模块

| 模块 | 方法 | 路径 | 鉴权 |
|------|------|------|------|
| file | POST | `/v2/file/upload` | 是 |
| file | GET | `/v2/file/list` | 是 |
| file | PUT | `/v2/file/:id` | 是 |
| file | DELETE | `/v2/file/:id` | 是 |
| example | GET | `/v2/example/info` | 否 |

## 中间件状态

全部 6 个中间件文件在生产代码中均已完全去全局化（零 `global` 导入）：

- `jwt.go` — `JWTAuthWithConfig(JWTConfig{ExpiresTime, BufferTime, SigningKey, BlacklistCheck})`
- `casbin_rbac.go` — `CasbinHandlerWithPrefix(prefix)`
- `cors.go` — `CorsByRules(config.CORS)`
- `error.go` — `GinRecovery(log, stack)`
- `operation.go` — `OperationRecord(db, log)`
- `limit_ip.go` — `NewLimiter(rdb, log, expire, limit)`

## 已删除

旧目录已全部移除：`api/`、`router/`、`service/`、`source/`、`core/`、`initialize/`

## 剩余工作（低优先级）

- `model/` 目录：25 个旧版 GORM 文件仍作为共享持久化层
- `utils/` 目录：`jwt.go`、`claims.go`、`casbin_util.go` 等仍引用 `global`，核心逻辑已迁入 `platform/auth`
- 测试文件在 `init()` 中设置 `global.GVA_*`（Go 常规模式）

## 如何添加新模块

```bash
cd server
go run ./cmd/modulegen -name <模块名>
```

然后在 `internal/modules/modules.go` 中注册。所有模块遵循依赖规则：

"platform 不得导入 modules。domain 不得导入框架。handler 不得触碰持久化。"
