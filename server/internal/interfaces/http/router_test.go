package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	global.GVA_LOG = zap.NewNop()
}

func TestRegisterV2Health(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterV2(engine)

	req := httptest.NewRequest(http.MethodGet, "/v2/health", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body response.Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Code != response.Success || body.Message != "ok" {
		t.Fatalf("response = %+v, want success ok", body)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T, want map[string]any", body.Data)
	}
	if data["status"] != "ok" {
		t.Fatalf("status = %v, want ok", data["status"])
	}
}

func TestRegisterV2Module(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterV2(engine, testModule{})

	req := httptest.NewRequest(http.MethodGet, "/v2/test-module/ping", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRegisterV2AuthenticatedModule(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterV2(engine, authenticatedTestModule{})

	req := httptest.NewRequest(http.MethodGet, "/v2/test-module/me", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

type testModule struct{}

func (testModule) RegisterHTTP(routes Routes) {
	routes.Public.GET("/test-module/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}

type authenticatedTestModule struct{}

func (authenticatedTestModule) RegisterHTTP(routes Routes) {
	routes.Authenticated.GET("/test-module/me", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
}
