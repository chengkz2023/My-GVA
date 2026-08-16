# 后端架构最终蓝图

> **历史文档**：本文档是模块化重构**之前**的规划稿，记录当时的设计意图与取舍过程。
> 部分内容（平台包清单、模块清单、目录约定）已与现状不一致。
> **当前权威文档以 `docs/backend-handoff.md`（现状报告）与 `server/docs/how-to-add-module.md`（开发指南）为准。**

## 1. 定位

本项目后端底座的目标不是复刻大厂平台工程，也不是继续沿用普通后台脚手架的快速 CRUD 结构，而是构建一套适合个人开发者长期复用的企业级模块化单体。

最终定位：

```text
个人开发友好的模块化单体
+ 清晰边界
+ 显式依赖
+ 可测试
+ 可迁移
+ 可扩展
+ 不过度设计
```

它应该吸收 Google、Uber、字节等大型科技公司的工程原则：

- 清晰简单。
- 依赖方向稳定。
- 避免可变全局状态。
- 平台能力统一治理。
- 模块边界可维护。
- 自动化防止架构腐化。

但不照搬大厂的组织形态：

- 不为几百人协作设计过重流程。
- 不一开始引入复杂平台体系。
- 不机械套完整 DDD。
- 不为了架构创建大量空目录。
- 不牺牲个人开发效率。

## 2. 当前问题

现有 gin-vue-admin 后端风格大致是：

```text
router -> api -> service -> model -> global.GVA_DB
```

这个结构适合快速生成管理后台，但作为长期项目底座存在几个问题：

- `global.GVA_DB`、`ServiceGroupApp`、`ApiGroupApp` 等全局对象让依赖不显式。
- handler/API 层和 service 层边界不够硬，业务容易散落。
- service 直接操作 GORM，复杂后会变成大函数仓库。
- request、response、database entity 容易混用。
- 模块之间缺少清晰边界，后期容易互相调用内部实现。
- 启动期自动迁移和种子数据逻辑对生产治理不够稳。

所以未来底座应当推倒重建后端架构，但可以保留有价值的能力：

- Gin。
- GORM。
- JWT。
- Casbin/RBAC 的核心能力。
- 用户、角色、菜单、权限、字典、文件等后台基础能力。
- 前端接口可阶段性兼容。

## 3. 核心原则

### 3.1 以模块为中心，而不是以技术层为中心

推荐：

```text
internal/modules/customer
internal/modules/order
internal/modules/system/user
internal/modules/system/role
```

避免所有业务都横向堆在：

```text
api/
service/
model/
router/
```

长期项目里，模块边界比技术层目录更重要。

### 3.2 简单模块用轻结构，复杂模块用完整分层

不要一刀切。

简单 CRUD 模块：

```text
customer/
  handler.go
  service.go
  repository.go
  model.go
  dto.go
  module.go
```

复杂业务模块：

```text
order/
  domain/
  application/
  infrastructure/mysql/
  transport/http/
  module.go
```

判断标准：

```text
后台配置表、简单数据维护 -> 轻结构
有状态流转 -> 完整分层
有金额、库存、审批、权限 -> 完整分层
会被多个模块复用 -> 完整分层
需要大量测试 -> 完整分层
```

### 3.3 显式依赖注入

业务代码不直接依赖：

```go
global.GVA_DB
global.GVA_CONFIG
global.GVA_LOG
```

推荐：

```go
type Service struct {
	repo Repository
	tx   transaction.Manager
}

func NewService(repo Repository, tx transaction.Manager) *Service {
	return &Service{repo: repo, tx: tx}
}
```

个人开发阶段优先手写构造函数。只有当装配明显复杂后，再考虑 Wire 或 Fx。

### 3.4 Handler 保持薄

HTTP handler 只负责：

- 参数绑定。
- 基础校验。
- 调用 service/usecase。
- 转换响应。

handler 不写复杂业务规则、不拼 SQL、不控制事务。

### 3.5 业务规则不要依赖框架

复杂模块的 `domain` 层禁止依赖：

- Gin。
- GORM。
- Zap。
- Viper。
- Redis。
- Casbin。
- HTTP DTO。

domain 只表达业务概念、状态流转、业务不变量和领域错误。

### 3.6 平台能力统一收敛

这些能力不应散落在业务模块里：

- 配置。
- 日志。
- 数据库连接。
- 事务。
- 认证。
- 权限。
- 错误码。
- 响应格式。
- 参数校验。
- 分页。
- 审计。

先实现常用能力，复杂能力预留扩展点，不急着全部落地。

### 3.7 生产环境不依赖 AutoMigrate

开发环境可以方便，生产环境必须可控。

推荐：

```text
开发环境：允许 AutoMigrate，可加载演示数据。
测试环境：使用 migration，可重建测试库。
生产环境：禁用 AutoMigrate，只允许版本化 migration。
```

### 3.8 不为架构而架构

底座的目的是让个人开发者长期更快，而不是显得更复杂。

禁止：

- 空洞的 DDD 包结构。
- 只有一个函数却拆五层。
- 为未来可能性过度抽象。
- 把每个 CRUD 都做成重型领域模型。
- 引入当前用不到的 MQ、Tracing、Outbox、CQRS。

## 4. 推荐目录结构

最终建议结构：

```text
server/
  cmd/
    admin-api/
      main.go

  internal/
    app/
      bootstrap/
        config.go
        logger.go
        database.go
        redis.go
        http.go
      container/
        container.go

    platform/
      auth/
      config/
      database/
      errors/
      logger/
      pagination/
      response/
      transaction/
      validator/

    modules/
      system/
        user/
        role/
        menu/
        permission/

      file/
      notification/
      tenant/

      business/
        example/

    interfaces/
      http/
        router.go
        middleware/

  api/
    openapi/
    proto/
    idl/

  configs/
    config.yaml
    config.example.yaml

  migrations/
    mysql/

  scripts/
  tests/
```

说明：

- `cmd/admin-api` 是启动入口。
- `internal/app` 负责启动、装配和生命周期。
- `internal/platform` 放通用基础能力。
- `internal/modules` 放业务模块。
- `internal/interfaces/http` 放 HTTP server、全局路由和中间件装配。
- `api` 放对外契约，例如 OpenAPI、Proto、IDL。
- `migrations` 管理数据库版本。

## 5. 模块结构规范

### 5.1 轻量模块

适合普通后台 CRUD。

```text
internal/modules/business/customer/
  handler.go
  service.go
  repository.go
  model.go
  dto.go
  module.go
```

职责：

- `handler.go`：HTTP 入参出参。
- `service.go`：业务逻辑和用例编排。
- `repository.go`：数据库访问。
- `model.go`：数据库 entity。
- `dto.go`：request/response DTO。
- `module.go`：模块注册。

约束：

- handler 不直接使用 GORM。
- repository 不写业务规则。
- service 不直接依赖 Gin。
- DTO 和数据库 entity 可以在简单场景下接近，但不要直接暴露敏感字段。

### 5.2 完整模块

适合复杂业务。

```text
internal/modules/business/order/
  domain/
    order.go
    status.go
    errors.go
    repository.go
    events.go

  application/
    service.go
    commands.go
    queries.go
    dto.go

  infrastructure/
    mysql/
      order_entity.go
      order_repository.go
      mapper.go

  transport/
    http/
      handler.go
      request.go
      response.go
      router.go

  module.go
```

职责：

- `domain`：业务实体、状态流转、核心规则、repository interface。
- `application`：用例编排、事务、权限、事件。
- `infrastructure/mysql`：GORM entity、查询、持久化实现。
- `transport/http`：Gin handler、HTTP DTO、路由。

## 6. 依赖规则

允许：

```text
handler -> service/application
service/application -> domain
service/application -> platform interface
repository implementation -> domain
repository implementation -> platform database
bootstrap/container -> all modules
modules -> platform
```

禁止：

```text
domain -> gin/gorm/zap/viper/redis/casbin
domain -> transport
domain -> infrastructure
handler -> GORM
handler -> repository implementation
platform -> modules
module A -> module B 的 repository implementation
module A -> module B 的 infrastructure
```

跨模块协作优先级：

1. 通过对方暴露的 application/service 接口。
2. 通过平台事件机制。
3. 通过明确的 facade。
4. 禁止直接访问对方内部 repository 或数据库表细节。

## 7. 平台层设计

### 7.1 config

封装配置读取。

业务代码不直接依赖 Viper。

推荐强类型配置：

```go
type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}
```

### 7.2 logger

封装日志能力。

业务层可以依赖轻量接口，避免到处 import zap。

### 7.3 database

负责：

- 创建数据库连接。
- 连接池配置。
- 健康检查。
- GORM 基础配置。
- 慢查询日志。

GORM 对象只应出现在 platform/database、repository、transaction manager。

### 7.4 transaction

统一事务入口。

推荐：

```go
err := tx.Within(ctx, func(ctx context.Context) error {
	return service.doSomething(ctx, cmd)
})
```

不要到处手写 `db.Transaction`。

### 7.5 errors

统一错误类型：

- validation
- unauthorized
- forbidden
- not found
- conflict
- internal

错误应能映射到：

- HTTP status。
- 业务 code。
- 用户可读 message。
- 日志等级。

### 7.6 response

统一响应格式，兼容前端。

业务层不依赖 response 包，只有 HTTP 层使用。

### 7.7 auth 与 permission

认证和权限分开：

```text
authentication: 当前用户是谁
authorization: 当前用户能否做某动作
data scope: 当前用户能看哪些数据
```

JWT、Casbin 可以保留，但封装为平台能力。

### 7.8 audit

审计日志可以先轻量实现。

只记录关键操作：

- 登录。
- 创建/修改/删除。
- 权限变更。
- 审批。
- 导出。

不要一开始做过重审计平台。

## 8. 数据库规范

### 8.1 Entity 与 DTO

简单模块可以让 entity 与 DTO 接近，但要避免直接暴露：

- 密码。
- Token。
- 内部状态。
- 删除标记。
- 审计字段。
- 敏感配置。

复杂模块必须拆分：

```text
domain model
database entity
request DTO
response DTO
```

### 8.2 Migration

目录：

```text
migrations/mysql/
  000001_create_users.up.sql
  000001_create_users.down.sql
```

原则：

- 生产只跑 migration。
- 每次结构变化有版本。
- migration 进入代码审查。
- seed 数据分类管理。

### 8.3 Seed 数据

分三类：

```text
required seed: 系统运行必需数据
demo seed: 演示数据，仅开发环境
tenant seed: 租户初始化数据
```

不要把演示数据和系统必需数据混在一起。

## 9. 权限模型

后台权限分为：

- 用户。
- 角色。
- 菜单。
- 按钮。
- API。
- 数据范围。

业务权限不应只依赖菜单/API 权限。

推荐业务代码这样表达权限：

```go
if err := authz.Can(ctx, actor, "order.approve", orderID); err != nil {
	return err
}
```

权限判断可以在 application 层做，handler 只负责解析身份。

## 10. 事件与异步

个人开发阶段不要一开始引入复杂 MQ。

推荐演进：

```text
阶段 1：同步调用
阶段 2：进程内 eventbus
阶段 3：outbox 表保证可靠投递
阶段 4：接入 MQ
```

先定义事件概念，等有真实需要再实现可靠事件。

## 11. 测试策略

### 11.1 轻量模块

优先测试 service。

repository 可用集成测试覆盖关键查询。

### 11.2 完整模块

测试分层：

```text
domain: 纯单元测试
application: mock/fake repository
infrastructure: 数据库集成测试
transport/http: handler 和响应测试
```

### 11.3 架构测试

后续加入 import 规则检查，例如：

- domain 不允许 import gin/gorm。
- handler 不允许 import infrastructure/mysql。
- platform 不允许 import modules。

## 12. 代码生成策略

个人开发需要效率，所以代码生成是底座能力之一。

允许生成：

- 轻量模块骨架。
- CRUD handler。
- CRUD service。
- repository 基础代码。
- DTO 模板。
- migration 模板。
- OpenAPI 文档。

禁止生成器覆盖：

- domain 规则。
- application 核心用例。
- 手写业务逻辑。
- 手写 mapper 中的复杂转换。

生成代码与手写代码要有清晰边界。

## 13. 与旧架构的迁移策略

### 阶段 1：建立新骨架

新增：

```text
cmd/admin-api
internal/app
internal/platform
internal/modules
internal/interfaces/http
configs
migrations
```

先让新入口能启动。

### 阶段 2：抽平台能力

优先抽：

- config。
- logger。
- database。
- response。
- errors。
- validator。
- transaction。
- auth。

### 阶段 3：迁移 system/user 作为范式

用户模块最适合做样板，因为它涉及：

- 登录。
- 当前用户。
- 用户列表。
- 创建用户。
- 修改密码。
- 角色关系。
- 权限上下文。

### 阶段 4：迁移权限体系

迁移：

- role。
- menu。
- permission。
- API permission。
- Casbin 封装。

### 阶段 5：迁移文件、通知、租户等通用模块

这些模块会成为未来项目复用资产。

### 阶段 6：新业务只走新架构

旧模块可以慢慢迁移，但新业务不再进入旧的 `api/service/model/router` 模式。

### 阶段 7：清理旧结构

当核心能力迁移完成后，逐步删除：

- 全局 service group。
- 全局 api group。
- 全局 router group。
- 业务对 `global.GVA_DB` 的直接依赖。
- 生产 AutoMigrate。
- 混杂的 request/response/entity。

## 14. 最小可行版本

第一版不要做太重。

MVP 目标：

```text
1. 新 cmd/admin-api 能启动。
2. config/logger/database/response/errors 可用。
3. transaction manager 可用。
4. HTTP router 可注册模块。
5. 一个轻量 CRUD 模块跑通。
6. 一个完整 user 模块跑通。
7. migration 工具链确定。
8. 新模块模板确定。
```

暂缓：

- OpenTelemetry。
- Outbox。
- MQ。
- 多数据库。
- 多租户复杂隔离。
- 完整代码生成器。
- Wire/Fx。
- 复杂 CQRS。

## 15. 最终验收标准

底座达到可长期复用时，应满足：

- 新增模块不需要修改全局变量。
- 新增模块有明确模板。
- 简单 CRUD 开发不比旧脚手架慢太多。
- 复杂业务可以清晰分层。
- handler 不写复杂业务。
- GORM 不进入 domain。
- 业务代码不直接依赖 `global.GVA_DB`。
- 生产数据库变更通过 migration。
- 权限、响应、错误、日志、事务有统一入口。
- 模块之间不会随意调用内部实现。
- 核心业务可以写单元测试。
- 架构规则可以逐步自动检查。

## 16. 一句话结论

这个底座的最佳形态不是大厂架构的复刻，也不是 CRUD 脚手架的延续，而是：

```text
为个人开发者优化的企业级模块化单体。
```

它应该足够先进，能支撑长期项目；也应该足够轻，不能拖慢个人开发速度。

## 17. 参考实践

- Go module organization: https://go.dev/doc/modules/layout
- Google Go Style Guide: https://google.github.io/styleguide/go/
- Uber Go Style Guide: https://github.com/uber-go/guide
- Google monorepo engineering paper: https://research.google/pubs/pub45424/
- CloudWeGo Hertz layout: https://www.cloudwego.io/docs/hertz/tutorials/toolkit/layout/
- CloudWeGo Kitex service governance: https://www.cloudwego.io/docs/kitex/tutorials/service-governance/
