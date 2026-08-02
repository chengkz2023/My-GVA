package example

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestModuleInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	apphttp.RegisterRoutes(engine, apphttp.Config{}, NewModule(nil))

	req := httptest.NewRequest(http.MethodGet, "/api/example/info", nil)
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
	if data["name"] != "business/example" {
		t.Fatalf("name = %v, want business/example", data["name"])
	}
}

func TestModuleMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	apphttp.RegisterRoutes(engine, apphttp.Config{}, NewModule(nil))

	req := httptest.NewRequest(http.MethodGet, "/api/example/missing", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	var body response.Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Code != response.Failure || body.Message != "example not found" {
		t.Fatalf("response = %+v, want failure example not found", body)
	}
}
