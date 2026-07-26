package bootstrap

import (
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/songzhibin97/gkit/cache/local_cache"
)

func InitBlackCache(jwt config.JWT) local_cache.Cache {
	dr, err := utils.ParseDuration(jwt.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(jwt.BufferTime)
	if err != nil {
		panic(err)
	}

	return local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
}
