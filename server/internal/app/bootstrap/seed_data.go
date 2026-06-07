package bootstrap

import "github.com/flipped-aurora/gin-vue-admin/server/global"

func EnsureSystemSeedData() {
	if global.GVA_DB == nil {
		return
	}
	// Seed data is now managed through migration files and the V2 admin API.
	// The AutoMigrate in RegisterTables() creates all required tables.
}
