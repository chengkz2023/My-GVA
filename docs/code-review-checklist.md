# 代码评审清单

评审对照项，分「必须」与「建议」。`go test ./internal -run TestArchitectureBoundaries` 已自动强制架构项，但评审仍需确认语义正确性。

## 架构（必须，多数由架构测试强制）

- [ ] `platform` 不 import `modules` 或顶层遗留包
- [ ] `domain` 框架无关（无 gin/gorm/zap/transport/infrastructure）
- [ ] `application` 无 gin（transport-free）
- [ ] `handler` 不直连持久化（无 gorm/infrastructure import）
- [ ] 模块不 import 其它模块的 `infrastructure`
- [ ] 无新增全局单例；依赖显式注入

## 安全（必须）

- [ ] 新写接口默认挂 `routes.Authenticated`；公开接口需说明理由并登记（如 `/example/greetings`、`/health`、`/login`、`/base/captcha`）
- [ ] 密码路径（创建/重置/修改）过 `platform/auth.PasswordPolicy`（默认 ≥8 位含字母数字）
- [ ] 登录类接口接防爆破：`platform/ratelimit` 按 IP + 账号双维度（参考 auth handler）
- [ ] SQL：参数化查询；排序等动态列名走白名单（先例：api 模块 `orderClause`）；无字符串拼接
- [ ] 错误输出不泄漏内部细节：Internal 错误文案不含底层 err；响应契约见 `docs/i18n.md`
- [ ] 日志不含请求体/敏感字段（操作审计只记元数据；panic 恢复只 dump 请求头）
- [ ] 鉴权路径符合 casbin 约定（策略路径不含 `/api` 前缀）；需要按角色过滤数据行的查询使用 `platform/dataauth.Scope`

## 质量（必须）

- [ ] 错误结构化：`apperrors.WithMessage(kind, ...)`，Kind 语义正确（Validation/Unauthorized/Forbidden/NotFound/Conflict/Internal）
- [ ] 分页走 `platform/pagination`（`Normalize` + `Result[T]`），页参数有默认值
- [ ] repo nil / `ErrRepositoryUnavailable` 优雅降级（列表返回空，不 panic）
- [ ] 返回结构向后兼容：新增字段只增不改；`code` 是契约（`docs/i18n.md`）
- [ ] application 新方法配单元测试；handler 新路由配测试
- [ ] 无 TODO/FIXME 遗留（当前仓库 TODO 计数 ≈ 0，评审应保持）

## 数据与迁移（必须）

- [ ] 系统表结构变更 → 属于脚手架版本发布，CHANGELOG 注明
- [ ] 业务表 → 版本化 SQL：`server/migrations/mysql/000N_xxx.sql`，禁止依赖 AutoMigrate（ADR-0003）
- [ ] 新 API 登记 `sys_apis`（seed 或同步）并为角色分配策略
- [ ] 种子数据幂等（`FirstOrCreate` / 按唯一键判断），老库升级能补种

## 前端（必须）

- [ ] 新增 UI 文案走 `t()`（vue-i18n），语言包在 `web/src/i18n/index.js`
- [ ] 请求逻辑判断用 `code` + HTTP 状态，不解析 `message`
- [ ] 401 由 `utils/request.js` 统一处理，页面不重复实现登出跳转
- [ ] 删除/危险操作有 `ElMessageBox` 确认
- [ ] 无 `e.srcElement`（用 `e.target`）；`window.open` 外链有 http/https 白名单

## 建议（可选）

- [ ] 大列表用 `KeepAlive` 缓存并正确清理事件监听（onUnmounted）
- [ ] 慢接口考虑 loading 状态与错误提示
- [ ] 常量/枚举值优先走字典模块（`/dictionary/types`）而非硬编码
