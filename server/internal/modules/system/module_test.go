package system

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	global.GVA_LOG = zap.NewNop()
}

func TestModuleRegistersChildren(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v2http.RegisterV2(engine, v2http.Config{}, NewModule(&container.Container{}))

	paths := map[string]int{
		"/v2/system/api/groups":       http.StatusUnauthorized,
		"/v2/system/api/list":         http.StatusUnauthorized,
		"/v2/system/api/all":          http.StatusUnauthorized,
		"/v2/system/api/policies/888": http.StatusUnauthorized,
		"/v2/system/auth/me":          http.StatusUnauthorized,
		"/v2/system/config/info":      http.StatusOK,
		"/v2/system/menu/tree":        http.StatusUnauthorized,
		"/v2/system/role/tree":        http.StatusUnauthorized,
		"/v2/system/status/info":      http.StatusOK,
		"/v2/system/user/list":        http.StatusUnauthorized,
		"/v2/system/user/me":          http.StatusUnauthorized,
		"/v2/system/version/info":     http.StatusOK,
	}
	for path, wantStatus := range paths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			engine.ServeHTTP(rec, req)
			if rec.Code != wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, wantStatus)
			}
		})
	}
}

func TestModuleRegistersProtectedPostRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v2http.RegisterV2(engine, v2http.Config{}, NewModule(&container.Container{}))

	req := httptest.NewRequest(http.MethodPost, "/v2/system/user/password", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestModuleRegistersProtectedPutRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v2http.RegisterV2(engine, v2http.Config{}, NewModule(&container.Container{}))

	req := httptest.NewRequest(http.MethodPut, "/v2/system/user/profile", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
