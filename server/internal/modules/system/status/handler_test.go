package status

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestModuleInfoWithoutDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/api")
	NewModule(&container.Container{}).RegisterHTTP(apphttp.Routes{Public: group, Authenticated: group})

	req := httptest.NewRequest(http.MethodGet, "/api/system/status/info", nil)
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
	checks, ok := data["checks"].(map[string]any)
	if !ok {
		t.Fatalf("checks type = %T, want map[string]any", data["checks"])
	}
	db, ok := checks["database"].(map[string]any)
	if !ok {
		t.Fatalf("database type = %T, want map[string]any", checks["database"])
	}
	if db["configured"] != false || db["ok"] != false {
		t.Fatalf("database = %v, want configured false ok false", db)
	}
}
