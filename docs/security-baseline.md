# 安全基线与加固指南

脚手架默认开启的安全能力，以及客户项目上线前的必做清单。定位：对外交付（见 `CONTEXT.md`「对外交付」与「安全基线」）。

## 默认开启（脚手架自带）

| 能力 | 位置 | 说明 |
|---|---|---|
| 登录防爆破限流 | `platform/ratelimit` + auth 模块 | IP + 账号双维度；失败计数、成功重置；超限 429；预算来自 `system.iplimit-count/iplimit-time` |
| 密码策略 | `platform/auth.PasswordPolicy` | 默认 ≥8 位且含字母数字；接入用户创建/修改密码/重置密码；管理员种子密码同样校验 |
| 管理员初始密码保护 | bootstrap 种子 | release 模式未设 `ADMIN_INITIAL_PASSWORD` 拒绝启动 |
| 图形验证码 | auth 模块 | 默认开启（`captcha.open-captcha: 1`），内存实现（单实例） |
| 操作审计 | `platform/audit.Sink` → `sys_operation_records` | 所有写操作（POST/PUT/DELETE/PATCH）自动记录；90 天自动清理 |
| JWT 黑名单 | `jwt_blacklists`（MySQL sha256） | 登出/改密后作废旧 token；本地缓存加速 |
| 超级管理员直通 | Casbin 中间件 | `authorityId == 888` 绕过策略检查——**生产必须为 888 账号设置强密码** |

## 生产上线必做清单

- [ ] **HTTPS 终结**：Nginx 配置 TLS（后端本身不处理 TLS），证书轮换纳入运维流程
- [ ] **更换 JWT 密钥**：`jwt.signing-key` 使用随机强密钥（≥32 字节），与开发环境不同
- [ ] **数据库凭据**：专用账号 + 最小权限，禁用 root；公网数据库必须限制来源 IP
- [ ] **`ADMIN_INITIAL_PASSWORD`**：release 模式显式设置强密码，且首次登录后通过「修改密码」更换
- [ ] **888 角色管控**：超级管理员账号只用于初始化与运维，日常操作使用普通角色
- [ ] **Redis 决策**：默认 `enable: false`（一切走 MySQL，功能完整）；若开启，Redis 必须设密码且后端与 Redis 均在内网
- [ ] **日志**：`zap.director` 指向专用目录，磁盘容量监控；日志不含请求体（审计只记元数据）
- [ ] **备份**：MySQL 定期全量 + binlog（备份策略由项目运维定，脚手架不含备份组件）

## 扩展接缝（客户项目按需替换）

| 接缝 | 接口 | 用法 |
|---|---|---|
| 密码策略 | `platform/auth.PasswordPolicy` | 等保复杂度/历史密码检查：实现 `Validate(password) error`，在 user 模块组装处注入 |
| 审计输出 | `platform/audit.Sink` | syslog/ELK：实现 `Record(ctx, Entry) error`，替换 bootstrap 的 `audit.NewMySQLSink` |
| 登录限流 | `platform/ratelimit.Limiter` | 多实例部署需换共享存储实现（Redis），接口与 auth handler 解耦 |
| 2FA / SSO / LDAP | 登录流程（auth 模块 `Login`） | 未内置，在 handler 的验证码校验与密码校验之间插入；接口占位见下方「改造点」 |

## 等保 2.0（二级）对照要点

| 要求 | 脚手架现状 |
|---|---|
| 身份鉴别 | 用户名+密码+验证码 ✅；2FA 需项目实现（扩展点） |
| 密码策略 | 默认 ≥8 位含字母数字 ✅（可按等保要求替换实现） |
| 登录失败处理 | 限流 + 失败计数 ✅（无账号锁定，需要时在 Policy/限流上扩展） |
| 审计 | 操作审计自动记录 ✅（输出可换 syslog） |
| 权限控制 | Casbin RBAC ✅；行级数据权限模式 `platform/dataauth` ✅ |
| 会话 | JWT + 黑名单 + 登出 ✅（无服务端会话，按项目评估） |
| 数据完整性/备份 | 由部署与运维保障（见生产必做清单） |

## 相关文档

- 错误码与 HTTP 状态契约：`docs/i18n.md`
- 配置项（密钥/验证码/限流参数）：`docs/configuration.md`
- 代码评审安全项：`docs/code-review-checklist.md`
