package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestModuleMeMissingActor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	c := &container.Container{
		Logger: zap.NewNop(),
	}
	apphttp.RegisterRoutes(engine, apphttp.Config{}, NewModule(c))
	req := httptest.NewRequest(http.MethodGet, "/api/system/auth/me", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
