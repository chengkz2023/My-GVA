package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestModuleMeMissingActor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	c := &container.Container{
		Logger: zap.NewNop(),
	}
	v2http.RegisterV2(engine, v2http.Config{}, NewModule(c))
	req := httptest.NewRequest(http.MethodGet, "/v2/system/auth/me", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
