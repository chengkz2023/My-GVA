# v1 → v2 差异文档

面向用过原版 gin-vue-admin（v1 时代的写法与习惯）的开发者；v2 即本脚手架 **v2.0.0**（DDD 模块化重构版）。
目标是避免按 v1 的肌肉记忆写出破坏 v2 架构或约定的代码。架构细节见 `server/docs/how-to-add-module.md` 与 `docs/backend-handoff.md`。

## 后端差异

| v1（旧习惯） | v2（现行约定） |
|---|---|
| 全局分层：`server/api`、`server/service`、`server/model`、`server/router`、`server/middleware`、`server/utils` | 模块化 DDD-lite：`internal/modules/<system\|business>/<name>/{domain,application,infrastructure,transport}` |
| 全局单例 `global.GVA_DB`、`global.GVA_LOG`、`global.GVA_VP` | 显式 DI：依赖由 `container.Container` 注入，零全局状态（架构测试强制） |
| 手写四层 + 在 router 集中注册 | `go run ./cmd/modulegen -name xxx` 生成骨架 + 照 `business/example` 写；注册在 `modules.go`（system 模块在 `system/module.go` 的 children） |
| `CasbinHandler()` + 888 硬编码跳过 | `CasbinHandlerWithPrefix(enforcer, prefix)` 参数化；超级管理员按 `authorityId == 888` 直通；**策略路径不含 `/api` 前缀** |
| 层级作用域检查只对非 888 角色生效 | 统一由 `use-strict-auth` 门控（默认 `false` = 不校验层级，见 `docs/adr` 与 role/user/menu 模块） |
| 888 默认放行写死在业务代码里 | 已移除所有硬编码 888 回退；只有 Casbin 中间件保留超级管理员直通 |
| 错误：handler 里直接 `response.Fail(code, msg)` | `platform/errors` 的 `Error{Kind, Code, Message}` + `response.Error(c, err)` 统一映射 HTTP 状态（契约见 `docs/i18n.md`） |
| 模型集中在 `model/` | 系统表在 `platform/database`；模块自己的表在模块 `infrastructure/mysql/models.go`（如 api、dictionary 模块） |
| 文件上传、任务调度、邮件、插件、swagger、断点续传、系统参数 | **已删除**（按项目自建，见 `docs/production-readiness.md`「已明确不做」） |
| 字典管理（v1 的 sysDictionaries） | **v2 内置重做**：`internal/modules/system/dictionary`（类型唯一、级联删除、`GET /dictionary/types` 业务引用接口 + 前端页面） |
| jwt 黑名单写 Redis | 写 MySQL `jwt_blacklists`（sha256 哈希，`char(64)` 唯一索引），本地缓存加速 |
| AutoMigrate 管所有表 | 系统表 AutoMigrate、**业务表版本化 SQL 迁移**（`migrations/mysql/0001_xxx.sql`，ADR-0003） |
| 角色树/菜单树直接查全表 | 有 strictAuth 门控与层级作用域；菜单树修复了 `parent_id = 0 OR IS NULL` 等历史问题 |

## 前端差异

- `view/superAdmin/*` 页面结构基本延续 v1，但**删除了一大批组件与依赖**：文件上传、`wangeditor`、`echarts`、`vue-office`、`vue-cropper`、二维码、拖拽排序、断点续传、`@form-create` 等（详见 CHANGELOG v1.0.0）。
- 字典管理页面：`view/superAdmin/dictionary/dictionary.vue`（v2 重做，接口与 v1 不同：`/dictionary/*`）。
- 新增页面：需要菜单种子（name + component 路径约定 `view/...`）+ 角色分配菜单；API 权限在「API 管理 → 同步 API」自动发现。
- **i18n 接缝**：vue-i18n v11 已接入，新增文案必须走 `t()`；存量硬编码中文不抽取（`docs/i18n.md`）。
- 请求封装：`utils/request.js` 401 自动登出；响应 `{code, message, msg, data}`，逻辑判断用 `code` 与 HTTP 状态，禁止解析 `message` 文案。

## 常见坑（按 v1 习惯最容易踩的）

1. 在 `platform` 层 import 任何模块 → 架构测试编译失败。
2. 在 `domain` 里 import gin/gorm → 架构测试失败。
3. 在 handler 里直接写 SQL/DB → 架构测试失败（handler 不能 import infrastructure/gorm）。
4. `sys_authority_menus.menuId` 是字符串列（v1 遗留），继续用字符串处理，别改成 uint。
5. 新 API 没进 `sys_apis`/没给角色分配策略 → 前端 403（888 角色除外）。
6. 密码相关功能绕过 `PasswordPolicy`（创建/重置/修改密码）→ 评审不过。
7. 新业务表图省事加进 AutoMigrate 列表 → 违反 ADR-0003，评审不过。
