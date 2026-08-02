package system

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func testContainer() *container.Container {
	return &container.Container{
		Logger: zap.NewNop(),
	}
}

func TestModuleRegistersChildren(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	apphttp.RegisterRoutes(engine, apphttp.Config{}, NewModule(testContainer()))

	paths := map[string]int{
		"/api/system/api/groups":       http.StatusUnauthorized,
		"/api/system/api/list":         http.StatusUnauthorized,
		"/api/system/api/all":          http.StatusUnauthorized,
		"/api/system/api/policies/888": http.StatusUnauthorized,
		"/api/system/auth/me":          http.StatusUnauthorized,
		"/api/system/config/info":      http.StatusOK,
		"/api/system/menu/tree":        http.StatusUnauthorized,
		"/api/system/role/tree":        http.StatusUnauthorized,
		"/api/system/status/info":      http.StatusOK,
		"/api/system/user/list":        http.StatusUnauthorized,
		"/api/system/user/me":          http.StatusUnauthorized,
		"/api/system/version/info":     http.StatusOK,
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
	apphttp.RegisterRoutes(engine, apphttp.Config{}, NewModule(testContainer()))

	req := httptest.NewRequest(http.MethodPost, "/api/system/user/password", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestModuleRegistersProtectedPutRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	apphttp.RegisterRoutes(engine, apphttp.Config{}, NewModule(testContainer()))

	req := httptest.NewRequest(http.MethodPut, "/api/system/user/profile", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
