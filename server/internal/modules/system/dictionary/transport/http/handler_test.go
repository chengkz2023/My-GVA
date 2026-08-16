package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/application"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestListDictionary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(platformauth.ContextWithActor(c.Request.Context(), platformauth.Actor{
			UserID: 1, AuthorityID: 888, Username: "admin", NickName: "Admin",
		}))
	})
	group := engine.Group("/api")
	NewHandler(application.NewService(&fakeRepository{})).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/api/dictionary/list?page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body response.Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T, want map[string]any", body.Data)
	}
	list, ok := data["list"].([]any)
	if !ok {
		t.Fatalf("list type = %T, want array", data["list"])
	}
	if len(list) != 0 {
		t.Fatalf("list = %v, want empty", list)
	}
}

func TestCreateDictionaryValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(platformauth.ContextWithActor(c.Request.Context(), platformauth.Actor{
			UserID: 1, AuthorityID: 888, Username: "admin", NickName: "Admin",
		}))
	})
	group := engine.Group("/api")
	NewHandler(application.NewService(&fakeRepository{})).Register(group)

	req := httptest.NewRequest(http.MethodPost, "/api/dictionary", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	var body response.Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if body.Code == response.Success {
		t.Fatalf("response = %+v, want error code", body)
	}
}
