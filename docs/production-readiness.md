# 生产成熟度验收清单（脚手架 v2.0.0）

来源：`/grill-with-docs` 会话结论（Q1–Q22）。P0 = 脚手架 v1.0 必须达成；P1 = 后续或随客户项目。

状态图例：✅ 已完成并提交 / ⬜ 未开始

## P0

1. ✅ **仓库卫生包**
   - 跟踪安全占位的 `config.yaml`（占位密码），README 快速开始改为 clone → 启动 真实可用
   - `Makefile` 修正模块路径与构建镜像（golang:1.24）
   - `web/Dockerfile` pnpm → npm，对齐 `package-lock.json`，内置 `nginx.conf`（同源代理）
   - 补 `server/.env.example`、`CHANGELOG.md` + 版本发布约定
   - 删除死掉的 CORS 配置段，文档化「同源代理」部署约定（`docs/deployment.md`）
2. ✅ **登录防爆破限流中间件**：`platform/ratelimit`（IP + 账号维度，登录失败计数、成功重置），接入 `iplimit-count`/`iplimit-time` 配置
3. ✅ **密码策略接缝**：`platform/auth.PasswordPolicy` + 默认实现（≥8 位含字母数字），接入用户创建 / 修改密码 / 重置密码三处
4. ✅ **管理员初始密码**：release 模式未设置 `ADMIN_INITIAL_PASSWORD` 拒绝启动；显式设置时过密码策略；debug 模式默认值仅告警
5. ✅ **验证码默认开启**：`config.example.yaml` / `config.yaml` 均 `open-captcha: 1`
6. ✅ **审计输出接缝**：`platform/audit.Sink` 接口 + MySQL 默认实现（operation-record 中间件改走接口）
7. ✅ **字典管理模块**：`system/dictionary` 全栈实现（DDD 四层 + 测试 + 种子示例字典 + 前端页面 `view/superAdmin/dictionary` + 菜单/API 增量种子）
8. ✅ **行级数据权限通用模式**：`platform/dataauth.Scope` 构造器 + example 模块 `ScopedList` 示范（含测试）
9. ✅ **i18n 接缝**：vue-i18n v11 就位、默认 zh-CN、字典页全部文案走 `t()` 作为示范、后端错误 code 契约文档（`docs/i18n.md`）
10. ✅ **文档包**：`docs/v1-v2-diff.md`、`docs/code-review-checklist.md`、`server/docs/how-to-add-module.md`（扩充平台能力/约定/清单）、`docs/deployment.md`（迁移与部署约定）、ADR-0001/0002/0003
11. ✅ **P1 顺带项**：JWT 黑名单本地缓存（`IsJwtBlacklistedCached`，只缓存肯定结果）、登录 cookie maxAge 读配置

## P1（后续或随客户项目）

- ⬜ 前端 Vitest 最小测试基建（当前前端零测试）
- ⬜ 参考 CI workflow 模板
- ⬜ 部署产物（compose/systemd/nginx 等）由实际项目决定，不入脚手架

## 已明确不做

- 文件上传/存储模块（按项目自建）
- 多租户（单租户边界，见 ADR-0001）
- 任务调度、通知、工作流（按项目自建）
- 存量中文文案全量 i18n 抽取（推迟到有海外交付的项目）

## 验证约定

每批改动：后端 `go build ./...` + `go test ./internal/...`，前端 `npm run build`，通过后立即提交（conventional commits）。

P0 全部完成并随 **v2.0.0** 发布，详见 `CHANGELOG.md` 的 `[v2.0.0]`。
