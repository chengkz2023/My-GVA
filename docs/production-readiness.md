# 生产成熟度验收清单（脚手架 v1.0）

来源：`/grill-with-docs` 会话结论（Q1–Q22）。P0 = 脚手架 v1.0 必须达成；P1 = 后续或随客户项目。

## P0

1. **仓库卫生包**
   - 跟踪安全占位的 `config.yaml`（占位密码），README 快速开始改为 clone → 复制配置 → 启动 真实可用
   - `Makefile` 修正模块路径（`flipped-aurora/gin-vue-admin` → 本项目）与构建镜像（golang:1.22 → 1.24）
   - `web/Dockerfile` pnpm → npm，对齐 `package-lock.json`
   - 补 `.env.example`、`CHANGELOG.md` + 版本发布约定
   - 删除死掉的 CORS 配置块，文档化「同源代理」部署约定
2. **登录防爆破限流中间件**（IP + 账号维度），接入现有 `iplimit-count`/`iplimit-time` 配置
3. **密码策略接缝**：`platform/auth` 策略接口 + 默认实现（≥8 位含字母数字），接入用户创建 / 修改密码 / 重置密码三处
4. **管理员初始密码**：release 模式未设置 `ADMIN_INITIAL_PASSWORD` 拒绝启动；debug 模式允许默认值并告警
5. **验证码默认开启**（示例配置 `open-captcha: 1`）
6. **审计输出接缝**：platform 层接口 + MySQL 默认实现（operation-record 改走接口）
7. **字典管理模块**：按 example 的 DDD 结构全新实现（后端 + 前端 + 测试 + 种子）
8. **行级数据权限通用模式**：platform 层 scope 构造器 + example 模块示范
9. **i18n 接缝**：vue-i18n 就位、默认 zh-CN、新增文案走 `t()`、后端错误 code 契约文档、示例文案
10. **文档包**：v1→v2 差异文档、代码评审清单、modulegen 标准流程、迁移约定（系统表 AutoMigrate / 业务表 SQL）、部署约定（同源代理）
11. **P1 顺带项**：JWT 黑名单缓存（`BlackCache` 接线或删除死代码）、登录 cookie maxAge 改为读配置

## P1（后续或随客户项目）

- 前端 Vitest 最小测试基建（当前前端零测试）
- 参考 CI workflow 模板
- 部署产物（compose/systemd/nginx 等）由实际项目决定，不入脚手架

## 已明确不做

- 文件上传/存储模块（按项目自建）
- 多租户（单租户边界，见 ADR）
- 任务调度、通知、工作流（按项目自建）
- 存量中文文案全量 i18n 抽取（推迟到有海外交付的项目）

## 验证约定

每批改动：后端 `go build ./...` + `go test ./internal/...`，前端 `npm run build`，通过后立即提交（conventional commits）。
