package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/business/example/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/business/example/infrastructure/memory"
	"github.com/gin-gonic/gin"
)

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	service := application.NewService(memory.NewRepository())
	NewHandler(service).Register(engine.Group("/api"))
	return engine
}

func TestListGreetings(t *testing.T) {
	engine := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/example/greetings", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "示例") {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestGetGreetingNotFound(t *testing.T) {
	engine := newTestRouter()
	req := httptest.NewRequest(http.MethodGet, "/api/example/greetings/999", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestCreateGreeting(t *testing.T) {
	engine := newTestRouter()
	body := `{"message":"hello","author":"tester"}`
	req := httptest.NewRequest(http.MethodPost, "/api/example/greetings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "hello") {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
