package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/domain"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

var testRoutes = func() gin.RoutesInfo { return nil }

func TestList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeRepository{}, nil, nil, false), testRoutes).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/v2/system/api/list?page=1&pageSize=10", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var body response.Body
	json.Unmarshal(rec.Body.Bytes(), &body)
	data, _ := body.Data.(map[string]any)
	if _, ok := data["list"]; !ok {
		t.Fatal("response missing 'list' field")
	}
}

func TestGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(platformauth.ContextWithActor(c.Request.Context(), platformauth.Actor{
			UserID:      1,
			AuthorityID: 888,
			Username:    "admin",
			NickName:    "Admin",
		}))
	})
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeRepository{}, nil, nil, false), testRoutes).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/v2/system/api/all", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestGroups(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeRepository{groups: []string{"admin", "api"}}, nil, nil, false), testRoutes).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/v2/system/api/groups", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestPolicies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	group := engine.Group("/v2")
	NewHandler(application.NewService(fakeRepository{}, &fakePolicyProvider{}, nil, false), testRoutes).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/v2/system/api/policies/888", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

type fakeRepository struct {
	groups []string
}

func (fakeRepository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.Api], error) {
	return pagination.Result[domain.Api]{
		List:     []domain.Api{{ID: 1, Path: "/api/test", Method: "GET"}},
		Total:    1,
		Page:     1,
		PageSize: 10,
	}, nil
}

func (fakeRepository) GetAll(ctx context.Context) ([]domain.Api, error) {
	return []domain.Api{{ID: 1, Path: "/api/test", Method: "GET"}}, nil
}

func (r fakeRepository) Groups(ctx context.Context) ([]string, error) {
	return r.groups, nil
}

type fakePolicyProvider struct{}

type fakePolicySyncer struct{}

func (fakePolicySyncer) RefreshPolicies() error                                    { return nil }
func (fakePolicySyncer) SyncPolicies(authorityID uint, policies []authz.Policy) error { return nil }

func (fakePolicyProvider) Policies(authorityID uint) ([]authz.Policy, error) {
	return []authz.Policy{{Path: "/api/test", Method: "GET"}}, nil
}
func (fakePolicySyncer) RemovePolicies(path, method string) error                     { return nil }
func (fakeRepository) Save(ctx context.Context, input domain.SaveApiInput) (uint, error) { return 1, nil }
func (fakeRepository) Delete(ctx context.Context, id uint) (domain.SaveApiInput, error) {
	return domain.SaveApiInput{Path: "/api/x", Method: "GET"}, nil
}
func (fakeRepository) DeleteByIds(ctx context.Context, ids []int) ([]domain.SaveApiInput, error) {
	return nil, nil
}
func (fakeRepository) FindByID(ctx context.Context, id uint) (domain.Api, error) {
	return domain.Api{ID: id, Path: "/api/test", Method: "GET"}, nil
}
func (fakeRepository) GetIgnored(ctx context.Context) ([]domain.Api, error)           { return nil, nil }
func (fakeRepository) Ignore(ctx context.Context, path, method string, flag bool) error { return nil }
func (fakeRepository) BatchCreate(ctx context.Context, apis []domain.SaveApiInput) error  { return nil }
func (fakeRepository) BatchDeleteByPathMethod(ctx context.Context, apis []domain.SaveApiInput) error {
	return nil
}
