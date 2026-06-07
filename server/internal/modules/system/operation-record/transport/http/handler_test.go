package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/operation-record/domain"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeRepo{})).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/v2/system/operation-record/list?page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var body response.Body
	json.Unmarshal(rec.Body.Bytes(), &body)
	if _, ok := body.Data.(map[string]any)["list"]; !ok {
		t.Fatal("response missing list field")
	}
}

func TestFindByIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeRepo{})).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/v2/system/operation-record/1", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestDeleteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeRepo{})).Register(group)

	req := httptest.NewRequest(http.MethodDelete, "/v2/system/operation-record/1", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

type fakeRepo struct{}

func (fakeRepo) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.Record], error) {
	return pagination.Result[domain.Record]{List: []domain.Record{{ID: 1, IP: "127.0.0.1"}}, Total: 1, Page: 1, PageSize: 10}, nil
}
func (fakeRepo) FindByID(ctx context.Context, id uint) (domain.Record, error) {
	return domain.Record{ID: id, IP: "127.0.0.1", Method: "GET", Path: "/api/test"}, nil
}
func (fakeRepo) Delete(ctx context.Context, id uint) error { return nil }
func (fakeRepo) DeleteByIds(ctx context.Context, ids []int) error { return nil }
