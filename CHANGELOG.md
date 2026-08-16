# Changelog

本文件是「模板复制」维护模式（docs/adr/0002）下，客户项目判断自己基于哪个脚手架版本的唯一依据。
脚手架每次版本发布必须更新本文件。

约定：
- 版本号 = 脚手架基版本（SemVer 语义：破坏性/新能力/修复）。
- 客户项目复制时，在标题下记录「Forked from vX.Y.Z」。
- 提交信息遵循 Conventional Commits。

## [Unreleased]

## [v2.0.0]

DDD 模块化重构 + P0 生产成熟度硬化后的第一个正式基版本。
分叉自上游 gin-vue-admin v2.8.8；上游历史标签（v0.9.0 / v1.0.0 / v2.x 系列）已从本仓库清理，
如需参考可重新 fetch 上游仓库（flipped-aurora/gin-vue-admin）。

### 架构

- DDD-lite 模块化单体：domain/application/infrastructure/transport 四层 + 显式 DI + 架构边界测试（6 条规则）
- 系统模块：用户/角色/菜单/API/字典/操作记录/配置/状态/版本
- 业务示例模块（example，新模块模板 + 行级数据权限示范）+ modulegen 生成器
- 前端：Vue3 + Vite + Element-Plus，动态路由 + 权限守卫；清理未用依赖与组件

### 安全基线（P0）

- 登录防爆破限流（platform/ratelimit，IP + 账号维度）
- 密码策略接缝（platform/auth.PasswordPolicy，默认 ≥8 位含字母数字）
- release 模式强制 `ADMIN_INITIAL_PASSWORD`；验证码默认开启
- 审计输出接缝（platform/audit.Sink，默认 MySQL）；JWT 黑名单 MySQL sha256 + 本地缓存
- 角色层级作用域统一由 `use-strict-auth` 门控
- npm 依赖零审计漏洞；删除 Vue CLI 遗留 devDeps（601 个包）

### 新能力

- 字典管理模块（后端 DDD 四层 + 前端页面 + 种子示例）
- 行级数据权限通用模式（platform/dataauth.Scope）
- i18n 接缝（vue-i18n v11 + 错误码契约）

### 修复（相对上游）

- 角色树 NULL 父节点、sys_ignore_apis 迁移、软删除、菜单树过滤、日志 core 并发竞争
- 前端 tabs `e.target`、warningBar 协议白名单、reload 无历史兜底、滚动条/焦点 a11y
- 仓库卫生：跟踪 config.yaml、.env.example、package-lock.json；Dockerfile 对齐 npm；删除死 CORS 配置
