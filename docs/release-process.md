# 发布流程（Release Process）

脚手架版本发布的固定步骤（客户项目复制基线的唯一依据，见 ADR-0002 与 `CHANGELOG.md`）。

## 何时发布

- P0/P1 类跨批功能完成，或累积了值得让客户项目回流的修复时。
- 版本号：SemVer——破坏性（major）/ 新能力（minor）/ 修复（patch）。

## 步骤

1. **改版本号**（两处）
   - `server/internal/platform/buildinfo/buildinfo.go` → `Version = "vX.Y.Z"`
   - `web/` 下 `npm version X.Y.Z --no-git-tag-version`（更新 package.json + package-lock.json）
2. **更新 CHANGELOG**
   - `## [Unreleased]` 内容移入 `## [vX.Y.Z]`，`[Unreleased]` 置空
   - 内容按「架构 / 安全基线 / 新能力 / 修复」归类（参照 v2.0.0 的写法）
3. **双端验证**
   - 后端：`cd server && go build ./... && go test ./internal/...`
   - 前端：`cd web && npm run build`
4. **提交**：`chore: bump version to vX.Y.Z — <一句话摘要>`
5. **打标签并推送**
   ```bash
   git tag -a vX.Y.Z -m "BoyKing Admin vX.Y.Z — <摘要>"
   git push origin vX.Y.Z
   git push origin main
   ```
6. **同步内网镜像**：将源头仓库的新 tag 同步到公司内网私有仓库（镜像同步方式由公司 Git 平台决定），并通知团队。

## 客户项目侧

- 评估回流：对照 CHANGELOG 决定 cherry-pick 哪些提交（ADR-0002）。
- 回流后记录：客户项目 CHANGELOG 中追加「Cherry-picked from vX.Y.Z」。
