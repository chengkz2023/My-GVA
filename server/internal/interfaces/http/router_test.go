package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestRegisterRoutesHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterRoutes(engine, Config{})

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
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

func TestRegisterRoutesModule(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterRoutes(engine, Config{}, testModule{})

	req := httptest.NewRequest(http.MethodGet, "/api/test-module/ping", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRegisterRoutesAuthenticatedModule(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	RegisterRoutes(engine, Config{}, authenticatedTestModule{})

	req := httptest.NewRequest(http.MethodGet, "/api/test-module/me", nil)
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
