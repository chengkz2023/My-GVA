package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/menu/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/menu/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
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
	group := engine.Group("/api")
	NewHandler(application.NewService(fakeRepository{}, nil, false)).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/api/system/menu/tree", nil)
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
	list, ok := data["menus"].([]any)
	if !ok || len(list) != 1 {
		t.Fatalf("menus = %v, want one item", data["menus"])
	}
}

type fakeRepository struct{}

func (fakeRepository) All(ctx context.Context) ([]domain.Menu, error) { return nil, nil }
func (fakeRepository) TreeByAuthority(ctx context.Context, authorityID uint) ([]domain.Menu, error) {
	return []domain.Menu{{
		ID:        1,
		ParentID:  0,
		Path:      "/dashboard",
		Name:      "superAdmin",
		Component: "dashboard/index",
		Sort:      1,
		Title:     "Dashboard",
		Icon:      "home",
	}}, nil
}

func (fakeRepository) Save(ctx context.Context, input domain.SaveMenuInput) (uint, error) { return 0, nil }
func (fakeRepository) Delete(ctx context.Context, id uint) error { return nil }
func (fakeRepository) FindByID(ctx context.Context, id uint) (domain.MenuDetail, error) { return domain.MenuDetail{}, nil }
func (fakeRepository) AssignMenus(ctx context.Context, authorityID uint, menuIDs []uint) error {
	return nil
}
