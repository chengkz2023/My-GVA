# 模板复制操作手册（Forking Guide）

客户项目从脚手架复制出独立副本（见 `CONTEXT.md`「模板复制」与 `docs/adr/0002-template-copy-maintenance.md`）。本手册给出复制后的全部改名步骤。

## 分发方式

- **源头**：个人 GitHub 仓库（当前 `chengkz2023/My-GVA`，tag `v2.0.0`）。
- **内网镜像**：公司内网私有仓库维护一份镜像，客户项目从镜像 fork/克隆。
- 复制后第一件事：在客户项目的 `CHANGELOG.md` 标题下记录 `Forked from v2.0.0`。

## 1. 复制仓库

```bash
# 方式 A：git clone 后删除 .git 重建（全新仓库）
git clone --depth 1 --branch v2.0.0 <镜像地址> my-project
cd my-project && rm -rf .git && git init

# 方式 B：在 GitLab/GitHub 上 fork，再 clone
```

## 2. 改名 Go 模块路径（必做）

脚手架的 Go module path `github.com/chengkz2023/My-GVA/server` 硬编码在 go.mod、全部 import、Makefile、Dockerfile 中。客户项目应改成自己的路径（如 `github.com/<公司>/<项目>/server`）：

```bash
# Windows PowerShell（在仓库根目录执行）
$old = "github.com/chengkz2023/My-GVA/server"
$new = "github.com/your-org/your-project/server"
Get-ChildItem -Recurse -Include *.go,go.mod,Makefile,Dockerfile -File |
  ForEach-Object { (Get-Content $_.FullName -Raw) -replace [regex]::Escape($old), $new | Set-Content -NoNewline $_.FullName }

# Linux / macOS（bash）
grep -rl "github.com/chengkz2023/My-GVA/server" --include="*.go" --include="go.mod" --include="Makefile" --include="Dockerfile" . \
  | xargs sed -i 's#github.com/chengkz2023/My-GVA/server#github.com/your-org/your-project/server#g'
```

改名后 `cd server && go build ./... && go test ./internal/...` 验证。

> 注意：`server/docs/how-to-add-module.md` 等文档里的示例 import 也引用旧路径（.md 共 19 处），建议一并替换或接受文档示例沿用脚手架路径。

## 3. 改名前端品牌与包名

| 位置 | 内容 |
|---|---|
| `web/package.json` `name` | `boyking-admin` → 项目包名 |
| `web/src/core/config.js` | `appName: 'BoyKing Admin'` → 项目名 |
| `web/index.html` `<title>` | BoyKing Admin → 项目名 |
| `web/nginx.conf`、`web/Dockerfile` | 品牌注释与镜像名（可选） |

## 4. 配置与密钥

- 复制 `server/configs/config.yaml`，修改：`jwt.signing-key`（必改）、`mysql.*`（连接与密码）、`zap.director`（日志目录）。
- 生产（release 模式）必须设置环境变量 `ADMIN_INITIAL_PASSWORD`，否则拒绝启动。
- 全部配置项含义见 `docs/configuration.md`。

## 5. 验证清单

```bash
cd server && go build ./... && go test ./internal/...   # 后端
cd web && npm install && npm run build                   # 前端
```

- [ ] 模块路径已替换且构建通过
- [ ] 前端 appName/标题已改
- [ ] `config.yaml` 密钥与数据库凭据已改
- [ ] CHANGELOG 已记录 `Forked from v2.0.0`
- [ ] （可选）重跑 agent-skills 配置以匹配新的 issue tracker（见 `docs/agents/issue-tracker.md`）
