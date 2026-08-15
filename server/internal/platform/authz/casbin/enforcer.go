package casbin

import (
	"errors"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const casbinModelText = `
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

// NewEnforcer 构造一个基于 DB adapter 的 casbin enforcer。
// 失败时返回错误，由调用方决定如何处理（不再使用包级全局单例）。
func NewEnforcer(db *gorm.DB) (*casbin.SyncedCachedEnforcer, error) {
	if db == nil {
		return nil, errors.New("casbin adapter: db is nil")
	}

	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromString(casbinModelText)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewSyncedCachedEnforcer(m, a)
	if err != nil {
		return nil, err
	}
	e.SetExpireTime(60 * 60)
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}
	return e, nil
}
