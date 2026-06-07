package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/application"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/domain"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestTree(t *testing.T) {
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
	NewHandler(application.NewService(fakeRepository{}, nil, nil, false)).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/v2/system/role/tree", nil)
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
	if !ok || len(list) != 1 {
		t.Fatalf("list = %v, want one item", data["list"])
	}
}

type fakeRepository struct{}

func (fakeRepository) Tree(ctx context.Context, authorityID uint, strict bool) ([]domain.Role, error) {
	return []domain.Role{{
		AuthorityID:   authorityID,
		AuthorityName: "admin",
		DefaultRouter: "dashboard",
	}}, nil
}

func (fakeRepository) Save(ctx context.Context, input domain.SaveRoleInput) error { return nil }
func (fakeRepository) Delete(ctx context.Context, id uint) (domain.SaveRoleInput, error) {
	return domain.SaveRoleInput{}, nil
}
func (fakeRepository) FindByID(ctx context.Context, id uint) (domain.Role, error) {
	return domain.Role{AuthorityID: id, AuthorityName: "test"}, nil
}
func (fakeRepository) FindMenuIDs(ctx context.Context, id uint) ([]uint, error) { return nil, nil }
func (fakeRepository) CopyMenusAndButtons(ctx context.Context, oldID, newID uint) error { return nil }
func (fakeRepository) SetDataAuthority(ctx context.Context, input domain.DataAuthorityInput) error { return nil }
func (fakeRepository) GetDataAuthorities(ctx context.Context, id uint) ([]uint, error) { return nil, nil }

func (fakeRepository) GetDescendantIDs(ctx context.Context, authorityID uint) ([]uint, error) {
	return nil, nil
}
func (fakeRepository) CheckAuthorityAuth(ctx context.Context, adminID, targetID uint) error {
	return nil
}
