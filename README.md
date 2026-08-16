# BoyKing Admin

基于 `gin + vue3 + element-plus` 的模块化后台管理脚手架（公司新项目标准起点，见 `CONTEXT.md`）。

## 快速开始

### 后端

```bash
cd server
# 配置：server/configs/config.yaml 已带安全占位默认值（clone 后可直接启动），
# 按需修改 mysql 连接；完整选项见 server/configs/config.example.yaml。
# 依赖：MySQL（必须，首次启动自动建表+种子数据）；Redis 默认关闭（可选）。
# 管理员初始密码：debug 模式默认 123456（告警）；release 模式必须设置环境变量 ADMIN_INITIAL_PASSWORD。
go run ./cmd/admin-api    # 默认 :8888
```

### 前端

```bash
cd web
npm install
npm run dev               # 默认 :8080，API 代理到 :8888
```

## 后端架构

DDD-lite 模块化单体，显式依赖注入，框架无关的领域层。

```
server/internal/
├── app/bootstrap/      启动编排（初始化、路由、数据库、定时器）
├── app/container/      DI 容器（Config, Logger, DB, Tx, Authorizer）
├── interfaces/http/    HTTP 路由 + 中间件
├── modules/system/     系统模块（auth/user/role/menu/api/operation-record/config/status/version）
├── modules/business/   业务模块（example — 开发示例，非真实业务）
└── platform/           平台共享层（auth/authz/casbin/config/database/errors/logger/...）
```

- 模块注册：`server/internal/modules/modules.go`
- JWT 核心逻辑（Claims、Token 生成/解析）位于 `platform/auth`，零全局依赖
- Casbin RBAC 认证与鉴权
- 中间件全部参数化：`JWTAuthWithConfig` / `CasbinHandler` / `GinRecovery(log)`
- 6 条架构边界规则自动检查：`go test ./internal -run TestArchitectureBoundaries`

## 工具

| 命令 | 用途 |
|------|------|
| `go run ./cmd/modulegen -name <name>` | 生成新模块代码 |

## 部署

同源代理约定与参考形态见 `docs/deployment.md`；前端 Nginx 示例 `web/nginx.conf`。

## 文档

- 词汇表：`CONTEXT.md`
- 决策记录：`docs/adr/`
- 生产成熟度验收清单：`docs/production-readiness.md`
- v1 → v2 差异（老团队必读）：`docs/v1-v2-diff.md`
- 代码评审清单：`docs/code-review-checklist.md`
- i18n 接缝与错误码契约：`docs/i18n.md`
- 部署约定：`docs/deployment.md`
- 新增模块指南：`server/docs/how-to-add-module.md`
- 详细 API 列表与架构说明：`docs/backend-handoff.md`
- 架构蓝图：`docs/backend-architecture-blueprint.md`
