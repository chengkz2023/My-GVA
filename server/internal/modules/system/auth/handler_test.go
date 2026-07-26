package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/gin-gonic/gin"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go.uber.org/zap"
)

func init() {
	global.GVA_LOG = zap.NewNop()
	global.GVA_CONFIG.JWT.SigningKey = "abcdefgh"
	global.GVA_CONFIG.JWT.ExpiresTime = "1h"
	global.GVA_CONFIG.JWT.BufferTime = "30m"
	global.BlackCache = local_cache.NewCache(local_cache.SetDefaultExpire(600))
}

func TestModuleMeMissingActor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v2http.RegisterV2(engine, v2http.Config{}, NewModule(&container.Container{}))
	req := httptest.NewRequest(http.MethodGet, "/v2/system/auth/me", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
