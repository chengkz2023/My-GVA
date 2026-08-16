# Changelog

本文件是「模板复制」维护模式（docs/adr/0002）下，客户项目判断自己基于哪个脚手架版本的唯一依据。
脚手架每次版本发布必须更新本文件。

约定：
- 版本号 = 脚手架基版本（SemVer 语义：破坏性/新能力/修复）。
- 客户项目复制时，在标题下记录「Forked from vX.Y.Z」。
- 提交信息遵循 Conventional Commits。

## [Unreleased]

### P0 生产成熟度硬化（docs/production-readiness.md）

- 仓库卫生：跟踪安全默认 `config.yaml`、`server/.env.example`、CHANGELOG、Makefile 修正、web Dockerfile 改 npm、删除死 CORS 配置、同源代理部署约定
- 安全基线：登录防爆破限流（platform/ratelimit）、密码策略接缝（platform/auth.PasswordPolicy）、release 模式强制 `ADMIN_INITIAL_PASSWORD`、验证码默认开启、审计输出接缝（platform/audit.Sink）
- JWT 黑名单本地缓存；登录 cookie maxAge 读配置
- 新增字典管理模块（system/dictionary，后端四层 + 前端页面 + 种子示例）
- 行级数据权限通用模式（platform/dataauth.Scope + example 示范）
- i18n 接缝：vue-i18n v11 + 字典页示范 + 错误码契约文档
- 文档包：v1→v2 差异、代码评审清单、模块开发指南扩充、3 个 ADR、CONTEXT 词汇表
- 依赖：移除 306 个残留包（npm install 修剪）；新增 vue-i18n

## [v1.0.0]

模块化重构后的第一个脚手架基版本：

- DDD-lite 模块化单体：domain/application/infrastructure/transport 四层 + 显式 DI + 架构边界测试
- JWT 认证（滑动续期）+ Casbin RBAC + 操作审计 + 登录登出黑名单
- 系统模块：用户/角色/菜单/API/配置/状态/版本/操作记录
- 业务示例模块（example，可作新模块模板）+ modulegen 生成器
- 前端：Vue3 + Vite + Element-Plus，动态路由 + 权限守卫，清理 25 个未用依赖
- 角色层级作用域检查统一由 use-strict-auth 门控
- 修复：角色树 NULL 父节点、sys_ignore_apis 迁移、软删除、菜单树过滤、日志 core 并发竞争等
