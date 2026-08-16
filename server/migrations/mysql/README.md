# MySQL Migrations（业务表）

业务表结构变更一律使用**版本化 SQL 迁移文件**，由客户项目自行评审与执行（见 `docs/adr/0003-migration-split.md`）。
脚手架系统表（用户/角色/菜单/API/字典/审计/黑名单等）由 GORM AutoMigrate 维护，**不属于本目录**。

## Migration files

```text
0001_<name>.up.sql
0001_<name>.down.sql
0002_<name>.up.sql
0002_<name>.down.sql
...
```

序号全局递增，一次变更一个序号；`up` 与 `down` 成对出现。

## Rules

- 每个 `up` 文件必须有对应的 `down` 文件。
- 生产环境的上线流程必须包含迁移文件的评审与执行步骤；不依赖 AutoMigrate 建业务表。
- 迁移文件内只放 DDL/DML；不要放环境相关值或密钥。
- 种子类数据与演示数据分离，演示数据只出现在脚手架系统表的种子逻辑（AutoMigrate + bootstrap seed），不进业务迁移。
