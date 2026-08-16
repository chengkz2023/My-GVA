# 配置与环境变量参考

配置文件：`server/configs/config.yaml`（默认安全占位，可直接启动）；全量示例与注释：`server/configs/config.example.yaml`。
加载顺序：命令行 `-c` > 环境变量 `GVA_CONFIG` > Gin 模式匹配（`config.debug.yaml` / `config.release.yaml` / `config.test.yaml`），`configs/` 子目录优先，最后回退 `config.yaml`。

## 环境变量

| 变量 | 必填 | 说明 |
|---|---|---|
| `ADMIN_INITIAL_PASSWORD` | release 模式必填 | 首次启动种子管理员（admin）的初始密码；须满足密码策略（≥8 位含字母数字）；debug 模式缺省时用 `123456` 并告警 |
| `GVA_CONFIG` | 否 | 指定配置文件路径（替代自动匹配） |

## system

| key | 默认 | 说明 | 生产建议 |
|---|---|---|---|
| `router-prefix` | `""` | 全局路由前缀（当前路由前缀由代码中 `/api` 常量控制，此项为预留） | 保持空 |
| `addr` | `8888` | HTTP 监听端口 | 按部署规划；对外只经 Nginx 反代 |
| `iplimit-count` | `15000` | 登录限流：窗口内每 IP/账号最大失败次数 | 按安全要求收紧（如 20） |
| `iplimit-time` | `3600` | 登录限流窗口（秒） | 与 count 配合 |
| `use-strict-auth` | `false` | 角色层级作用域校验开关（见 CONTEXT「安全基线」；true 时角色/用户/菜单操作只允许作用域内角色） | 多级管理员时开启 |
| `disable-auto-migrate` | `false` | 关闭启动时的系统表 AutoMigrate | 保持 false（系统表由脚手架维护，ADR-0003） |

## jwt

| key | 默认 | 说明 | 生产建议 |
|---|---|---|---|
| `signing-key` | 占位串 | HS256 签名密钥 | **必改**：随机 ≥32 字节，与开发环境不同 |
| `expires-time` | `7d` | token 有效期（也决定登录 cookie maxAge） | 按安全要求（如 2h–1d） |
| `buffer-time` | `1d` | 滑动续期窗口 | 小于 expires-time |
| `issuer` | `boyking-admin` | JWT 签发者 | 可改项目名 |

## zap（日志）

| key | 默认 | 说明 |
|---|---|---|
| `level` | `info` | 日志级别（debug/info/warn/error…） |
| `prefix` | `[boyking-admin/server]` | 每行时间前缀 |
| `format` | `console` | `console` 或 `json` |
| `director` | `log` | 日志目录（相对 server 工作目录），按 `级别/日期/级别.log` 切分 |
| `encode-level` | `LowercaseColorLevelEncoder` | 级别编码风格 |
| `stacktrace-key` | `stacktrace` | 堆栈字段名 |
| `show-line` | `true` | 是否打印调用位置 |
| `log-in-console` | `true` | 同时输出到控制台 |
| `retention-day` | `-1` | 日志保留天数（-1 = 永不清理；正整数按目录 mtime 清理） |

运维要点：磁盘容量监控针对 `director` 目录；生产建议 `retention-day` 设为 30 或 90。

## captcha（验证码）

| key | 默认 | 说明 |
|---|---|---|
| `key-long` | `6` | 验证码长度 |
| `img-width` / `img-height` | `240` / `80` | 图片尺寸 |
| `open-captcha` | `1` | 验证码开关（0=关闭；>0 开启） |
| `open-captcha-timeout` | `3600` | 验证码缓存有效期（秒，内存实现） |

## mysql

| key | 默认 | 说明 | 生产建议 |
|---|---|---|---|
| `path` / `port` | `127.0.0.1` / `3306` | 连接地址 | 专用账号，禁 root |
| `db-name` | `boyking_admin` | 库名 | — |
| `username` / `password` | 占位 | 凭据 | **必改** |
| `config` | `charset=utf8mb4&parseTime=True&loc=Local` | DSN 参数 | 保持 |
| `max-idle-conns` | `10` | 空闲连接数 | 按负载 |
| `max-open-conns` | `100` | 最大连接数 | 按负载 |
| `log-mode` | `error` | GORM SQL 日志级别（silent/error/warn/info） | 生产用 silent/error |
| `log-zap` | `false` | SQL 日志走 zap（预留） | — |
| `prefix` | `""` | 表前缀（预留） | 空 |
| `singular` | `false` | 表名复数 | 保持 false |
| `engine` | `""` | 建表引擎（空 = InnoDB 默认） | 空 |

## redis（可选，默认关闭）

| key | 默认 | 说明 |
|---|---|---|
| `enable` | `false` | 是否启用。**关闭时功能完整**（黑名单走 MySQL、验证码走内存、无缓存依赖） |
| `addr` / `password` / `db` | `127.0.0.1:6379` / 空 / `0` | 连接；启用则必须设密码 |
| `useCluster` / `clusterAddrs` | `false` / `[]` | 集群模式 |

注意：`enable: true` 时 Redis 必须可达（初始化失败会直接终止启动）；多实例部署时登录限流与验证码需自行替换为共享存储实现（见 `docs/security-baseline.md`）。

## 备份与升级要点

- **备份**：MySQL 定期全量 dump + binlog；`sys_operation_records`/`jwt_blacklists` 可排除（可重建数据）；备份策略由项目运维制定。
- **升级（v2.x → v2.y）**：脚手架系统表随版本由 AutoMigrate 自动演进，种子幂等补种（新菜单/新 API/新字典自动补）；业务表按 `server/migrations/mysql/` 的 SQL 迁移执行；修复回流用 cherry-pick（ADR-0002），以 CHANGELOG 判断是否需要。
