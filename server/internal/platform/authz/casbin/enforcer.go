package casbin

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	enforcer *casbin.SyncedCachedEnforcer
	once     sync.Once
)

func GetEnforcer(db *gorm.DB) *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		if db == nil {
			zap.L().Error("casbin adapter: db is nil")
			return
		}
		a, err := gormadapter.NewAdapterByDB(db)
		if err != nil {
			zap.L().Error("casbin adapter: failed to create adapter", zap.Error(err))
			return
		}
		text := `
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		m, err := model.NewModelFromString(text)
		if err != nil {
			zap.L().Error("casbin model: failed to load model", zap.Error(err))
			return
		}
		enforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
		enforcer.SetExpireTime(60 * 60)
		_ = enforcer.LoadPolicy()
	})
	return enforcer
}
