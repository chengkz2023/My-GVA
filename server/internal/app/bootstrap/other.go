package bootstrap

import (
	"github.com/chengkz2023/My-GVA/server/config"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/songzhibin97/gkit/cache/local_cache"
)

func InitBlackCache(jwt config.JWT) local_cache.Cache {
	dr, err := platformauth.ParseDuration(jwt.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = platformauth.ParseDuration(jwt.BufferTime)
	if err != nil {
		panic(err)
	}

	return local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
}
