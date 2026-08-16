package bootstrap

import (
	"os"
	"strconv"

	apimysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/api/infrastructure/mysql"
	dictionarymysql "github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/infrastructure/mysql"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	adminAuthorityID = 888
	adminUsername    = "admin"
)

func EnsureSystemSeedData(db *gorm.DB, log *zap.Logger) {
	if db == nil {
		return
	}
	seedAuthority(db, log)
	seedAdminUser(db, log)
	seedMenus(db, log)
	seedAPIs(db, log)
	seedDictionaries(db, log)
	log.Info("seed data complete")
}

func seedAuthority(db *gorm.DB, log *zap.Logger) {
	var count int64
	db.Model(&platformdb.SysAuthority{}).Where("authority_id = ?", adminAuthorityID).Count(&count)
	if count > 0 {
		return
	}
	auth := platformdb.SysAuthority{
		AuthorityId:   adminAuthorityID,
		AuthorityName: "超级管理员",
		DefaultRouter: "authority",
	}
	if err := db.Create(&auth).Error; err != nil {
		log.Error("seed authority failed", zap.Error(err))
	}
}

func seedAdminUser(db *gorm.DB, log *zap.Logger) {
	var count int64
	db.Model(&platformdb.SysUser{}).Where("username = ?", adminUsername).Count(&count)
	if count > 0 {
		return
	}
	initialPassword := os.Getenv("ADMIN_INITIAL_PASSWORD")
	switch {
	case initialPassword != "":
		policyErr := platformauth.DefaultPasswordPolicy{}.Validate(initialPassword)
		if policyErr != nil {
			log.Fatal("ADMIN_INITIAL_PASSWORD does not satisfy the password policy", zap.Error(policyErr))
		}
	case gin.Mode() == gin.ReleaseMode:
		log.Fatal("ADMIN_INITIAL_PASSWORD is required in release mode — refusing to seed admin with a default password")
	default:
		initialPassword = "123456"
		log.Warn("ADMIN_INITIAL_PASSWORD not set, using default admin password '123456' — change it immediately")
	}
	hashedPwd, err := platformauth.NewBcryptPasswordHasher().Hash(initialPassword)
	if err != nil {
		log.Error("seed admin password hash failed", zap.Error(err))
		return
	}
	user := platformdb.SysUser{
		UUID:        uuid.New(),
		Username:    adminUsername,
		Password:    hashedPwd,
		NickName:    "管理员",
		HeaderImg:   "https://qmplusimg.henrongyi.top/gva_header.jpg",
		AuthorityId: adminAuthorityID,
		Enable:      1,
	}
	if err := db.Create(&user).Error; err != nil {
		log.Error("seed admin user failed", zap.Error(err))
		return
	}
	db.Create(&platformdb.SysUserAuthority{
		SysUserId:               user.ID,
		SysAuthorityAuthorityId: adminAuthorityID,
	})
}

func seedMenus(db *gorm.DB, log *zap.Logger) {
	var count int64
	db.Model(&platformdb.SysBaseMenu{}).Count(&count)
	if count > 0 {
		return
	}

	menus := []platformdb.SysBaseMenu{
		{
			ParentId: 0, Path: "admin", Name: "superAdmin", Hidden: false,
			Component: "view/superAdmin/index.vue", Sort: 1,
			Meta: platformdb.Meta{Title: "超级管理员", Icon: "user"},
		},
		{
			ParentId: 1, Path: "authority", Name: "authority", Hidden: false,
			Component: "view/superAdmin/authority/authority.vue", Sort: 1,
			Meta: platformdb.Meta{Title: "角色管理", Icon: "avatar"},
		},
		{
			ParentId: 1, Path: "menu", Name: "menu", Hidden: false,
			Component: "view/superAdmin/menu/menu.vue", Sort: 2,
			Meta: platformdb.Meta{Title: "菜单管理", Icon: "tickets", KeepAlive: true},
		},
		{
			ParentId: 1, Path: "api", Name: "api", Hidden: false,
			Component: "view/superAdmin/api/api.vue", Sort: 3,
			Meta: platformdb.Meta{Title: "API管理", Icon: "platform", KeepAlive: true},
		},
		{
			ParentId: 1, Path: "user", Name: "user", Hidden: false,
			Component: "view/superAdmin/user/user.vue", Sort: 4,
			Meta: platformdb.Meta{Title: "用户管理", Icon: "coordinate"},
		},
		{
			ParentId: 1, Path: "operation", Name: "operation", Hidden: false,
			Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 5,
			Meta: platformdb.Meta{Title: "操作历史", Icon: "pie-chart"},
		},
		{
			ParentId: 1, Path: "dictionary", Name: "dictionary", Hidden: false,
			Component: "view/superAdmin/dictionary/dictionary.vue", Sort: 6,
			Meta: platformdb.Meta{Title: "字典管理", Icon: "notebook"},
		},
	}
	for i := range menus {
		if err := db.Create(&menus[i]).Error; err != nil {
			log.Error("seed menu failed", zap.Error(err))
			continue
		}
		db.Create(&platformdb.SysAuthorityMenu{
			MenuId:      strconv.Itoa(int(menus[i].ID)),
			AuthorityId: strconv.Itoa(adminAuthorityID),
		})
	}
}

func seedAPIs(db *gorm.DB, log *zap.Logger) {
	var count int64
	db.Model(&apimysql.SysApi{}).Count(&count)
	if count > 0 {
		return
	}
	apis := []apimysql.SysApi{
		{Path: "/health", Description: "健康检查", ApiGroup: "public", Method: "GET"},
		{Path: "/login", Description: "登录", ApiGroup: "public", Method: "POST"},
		{Path: "/base/captcha", Description: "验证码", ApiGroup: "public", Method: "POST"},
		{Path: "/system/auth/me", Description: "当前用户信息", ApiGroup: "auth", Method: "GET"},
		{Path: "/system/user/list", Description: "用户列表", ApiGroup: "user", Method: "GET"},
		{Path: "/system/user/me", Description: "个人中心", ApiGroup: "user", Method: "GET"},
		{Path: "/system/user", Description: "创建用户", ApiGroup: "user", Method: "POST"},
		{Path: "/system/user/profile", Description: "更新资料", ApiGroup: "user", Method: "PUT"},
		{Path: "/system/user/password", Description: "修改密码", ApiGroup: "user", Method: "POST"},
		{Path: "/system/role/tree", Description: "角色树", ApiGroup: "role", Method: "GET"},
		{Path: "/system/role", Description: "创建角色", ApiGroup: "role", Method: "POST"},
		{Path: "/system/role/copy", Description: "复制角色", ApiGroup: "role", Method: "POST"},
		{Path: "/system/role/data-authority", Description: "数据权限", ApiGroup: "role", Method: "POST"},
		{Path: "/system/role/data-authorities", Description: "获取数据权限", ApiGroup: "role", Method: "GET"},
		{Path: "/system/menu/tree", Description: "菜单树", ApiGroup: "menu", Method: "GET"},
		{Path: "/system/menu", Description: "创建菜单", ApiGroup: "menu", Method: "POST"},
		{Path: "/system/menu/authority", Description: "分配菜单权限", ApiGroup: "menu", Method: "POST"},
		{Path: "/system/api/list", Description: "API列表", ApiGroup: "api", Method: "GET"},
		{Path: "/system/api/all", Description: "全部API", ApiGroup: "api", Method: "GET"},
		{Path: "/system/api/groups", Description: "API分组", ApiGroup: "api", Method: "GET"},
		{Path: "/system/api/policies", Description: "更新策略", ApiGroup: "api", Method: "PUT"},
		{Path: "/system/api/policies/:authorityId", Description: "角色策略", ApiGroup: "api", Method: "GET"},
		{Path: "/system/api", Description: "创建API", ApiGroup: "api", Method: "POST"},
		{Path: "/system/api/batch-delete", Description: "批量删除API", ApiGroup: "api", Method: "POST"},
		{Path: "/system/api/fresh-casbin", Description: "刷新Casbin", ApiGroup: "api", Method: "POST"},
		{Path: "/system/api/sync", Description: "同步API", ApiGroup: "api", Method: "GET"},
		{Path: "/system/api/ignore", Description: "忽略API", ApiGroup: "api", Method: "POST"},
		{Path: "/system/api/batch-sync", Description: "批量同步API", ApiGroup: "api", Method: "POST"},
		{Path: "/system/operation-record/list", Description: "操作记录列表", ApiGroup: "operation-record", Method: "GET"},
		{Path: "/system/operation-record/batch-delete", Description: "批量删除", ApiGroup: "operation-record", Method: "POST"},
		{Path: "/dictionary/list", Description: "字典列表", ApiGroup: "dictionary", Method: "GET"},
		{Path: "/dictionary/types", Description: "字典引用数据", ApiGroup: "dictionary", Method: "GET"},
		{Path: "/dictionary", Description: "创建字典", ApiGroup: "dictionary", Method: "POST"},
		{Path: "/dictionary/:id", Description: "更新字典", ApiGroup: "dictionary", Method: "PUT"},
		{Path: "/dictionary/:id", Description: "删除字典", ApiGroup: "dictionary", Method: "DELETE"},
		{Path: "/dictionary/details", Description: "字典项列表", ApiGroup: "dictionary", Method: "GET"},
		{Path: "/dictionary/details", Description: "创建字典项", ApiGroup: "dictionary", Method: "POST"},
		{Path: "/dictionary/details/:id", Description: "更新字典项", ApiGroup: "dictionary", Method: "PUT"},
		{Path: "/dictionary/details/:id", Description: "删除字典项", ApiGroup: "dictionary", Method: "DELETE"},
		{Path: "/system/config/info", Description: "系统配置", ApiGroup: "config", Method: "GET"},
		{Path: "/system/status/info", Description: "服务状态", ApiGroup: "status", Method: "GET"},
		{Path: "/system/version/info", Description: "版本信息", ApiGroup: "version", Method: "GET"},
	}
	for _, api := range apis {
		db.Where("path = ? AND method = ?", api.Path, api.Method).FirstOrCreate(&api)
	}
}

// seedDictionaries 种一个示例字典（性别），展示字典管理的完整用法。
func seedDictionaries(db *gorm.DB, log *zap.Logger) {
	var count int64
	db.Model(&dictionarymysql.SysDictionary{}).Where("type = ?", "gender").Count(&count)
	if count > 0 {
		return
	}
	dictionary := dictionarymysql.SysDictionary{Type: "gender", Name: "性别", Sort: 1, Status: 1}
	if err := db.Create(&dictionary).Error; err != nil {
		log.Error("seed dictionary failed", zap.Error(err))
		return
	}
	details := []dictionarymysql.SysDictionaryDetail{
		{DictionaryID: dictionary.ID, Label: "男", Value: "male", Sort: 1, Status: 1},
		{DictionaryID: dictionary.ID, Label: "女", Value: "female", Sort: 2, Status: 1},
		{DictionaryID: dictionary.ID, Label: "未知", Value: "unknown", Sort: 3, Status: 1},
	}
	for i := range details {
		if err := db.Create(&details[i]).Error; err != nil {
			log.Error("seed dictionary detail failed", zap.Error(err))
		}
	}
}
