package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/domain"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
	"github.com/gin-gonic/gin"
)

func TestListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeFileRepo{}, "/tmp")).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/v2/file/list?page=1&pageSize=10", nil)
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
	NewHandler(application.NewService(fakeFileRepo{files: []domain.File{{ID: 1, Key: "test.png"}}}, "/tmp")).Register(group)

	req := httptest.NewRequest(http.MethodDelete, "/v2/file/1", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestUpdateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeFileRepo{files: []domain.File{{ID: 1}}}, "/tmp")).Register(group)

	body, _ := json.Marshal(map[string]string{"name": "new.png", "tag": "png"})
	req := httptest.NewRequest(http.MethodPut, "/v2/file/1", bytesReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

type bytesReader []byte

func (b bytesReader) Read(p []byte) (int, error) { return copy(p, b), nil }
func (b bytesReader) Close() error               { return nil }

type fakeFileRepo struct {
	files []domain.File
}

func (fakeFileRepo) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.File], error) {
	return pagination.Result[domain.File]{List: []domain.File{{ID: 1, Name: "test.png"}}, Total: 1, Page: 1, PageSize: 10}, nil
}
func (fakeFileRepo) FindByID(ctx context.Context, id uint) (domain.File, error) {
	return domain.File{}, nil
}
func (r fakeFileRepo) Create(ctx context.Context, file domain.File) (domain.File, error) {
	return domain.File{ID: 1, Name: file.Name}, nil
}
func (r fakeFileRepo) Update(ctx context.Context, id uint, name string, tag string) (domain.File, error) {
	return domain.File{ID: id, Name: name, Tag: tag}, nil
}
func (r fakeFileRepo) Delete(ctx context.Context, id uint) (domain.File, error) {
	if len(r.files) > 0 {
		return r.files[0], nil
	}
	return domain.File{}, nil
}
