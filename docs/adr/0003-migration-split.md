# 迁移分工：系统表 AutoMigrate，业务表 SQL

脚手架系统表（用户/角色/菜单/API/操作审计/黑名单等）由 GORM AutoMigrate 在启动时维护，默认开启（含生产）；业务表一律使用版本化 SQL 迁移文件（`server/migrations/mysql/` 下 `0001_xxx.sql` 命名），由客户项目自行执行与评审。

**Considered Options**: 全部 SQL 迁移——被拒，系统表随脚手架版本演进频繁且由脚手架自身控制，AutoMigrate 保证新副本开箱即建表；全部 AutoMigrate——被拒，业务表结构变更需要可审计、可评审的变更记录（等保/客户审计要求），且 GORM 对复杂变更（改列类型、数据搬迁）不可靠。

**Consequences**: 业务表模型不得依赖 AutoMigrate 建表；客户项目上线流程必须包含 SQL 迁移文件的评审与执行步骤；脚手架对系统表的结构变更随版本发布并在 CHANGELOG 中注明。
