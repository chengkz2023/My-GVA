package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestModuleInfoRedactsSecrets(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Server{}
	cfg.System.Addr = 8888
	cfg.JWT.SigningKey = "secret"
	cfg.JWT.Issuer = "boyking-admin"
	cfg.Mysql.Path = "127.0.0.1"
	cfg.Mysql.Dbname = "boyking_admin"
	cfg.Mysql.Password = "mysql-secret"
	cfg.Redis.Password = "redis-secret"

	engine := gin.New()
	v2http.RegisterV2(engine, v2http.Config{}, NewModule(&container.Container{Config: cfg}))

	req := httptest.NewRequest(http.MethodGet, "/v2/system/config/info", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := rec.Body.String(); body == "" || containsAny(body, "secret", "mysql-secret", "redis-secret") {
		t.Fatalf("response leaked secret or was empty: %s", body)
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
	configData, ok := data["config"].(map[string]any)
	if !ok {
		t.Fatalf("config type = %T, want map[string]any", data["config"])
	}
	jwtData, ok := configData["jwt"].(map[string]any)
	if !ok || jwtData["signingKeySet"] != true {
		t.Fatalf("jwt data = %v, want signingKeySet true", configData["jwt"])
	}
	mysqlData, ok := configData["mysql"].(map[string]any)
	if !ok || mysqlData["passwordSet"] != true {
		t.Fatalf("mysql data = %v, want passwordSet true", configData["mysql"])
	}
	redisData, ok := configData["redis"].(map[string]any)
	if !ok || redisData["passwordSet"] != true {
		t.Fatalf("redis data = %v, want passwordSet true", configData["redis"])
	}
}

func containsAny(s string, values ...string) bool {
	for _, value := range values {
		if value != "" && strings.Contains(s, value) {
			return true
		}
	}
	return false
}
