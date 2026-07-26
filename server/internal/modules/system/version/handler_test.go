package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestModuleInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	v2http.RegisterV2(engine, v2http.Config{}, NewModule(nil))

	req := httptest.NewRequest(http.MethodGet, "/v2/system/version/info", nil)
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
	if data["version"] != global.Version {
		t.Fatalf("version = %v, want %s", data["version"], global.Version)
	}
	if data["appName"] != global.AppName {
		t.Fatalf("appName = %v, want %s", data["appName"], global.AppName)
	}
}
