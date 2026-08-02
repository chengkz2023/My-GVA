package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/user/application"
	"github.com/chengkz2023/My-GVA/server/internal/modules/system/user/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
)

func TestMe(t *testing.T) {
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
	NewHandler(application.NewService(nil, platformauth.NewBcryptPasswordHasher(), nil)).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/api/system/user/me", nil)
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
	user, ok := data["user"].(map[string]any)
	if !ok {
		t.Fatalf("user type = %T, want map[string]any", data["user"])
	}
	if user["username"] != "admin" || user["source"] != "token" {
		t.Fatalf("user = %v, want admin token source", user)
	}
}

func TestList(t *testing.T) {
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
	NewHandler(application.NewService(nil, platformauth.NewBcryptPasswordHasher(), nil)).Register(group)

	req := httptest.NewRequest(http.MethodGet, "/api/system/user/list?page=0&pageSize=200&username=admin", nil)
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
	if data["page"] != float64(1) || data["pageSize"] != float64(100) {
		t.Fatalf("pagination = %v/%v, want normalized 1/100", data["page"], data["pageSize"])
	}
}

func TestChangePassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeRepository{passwordHash: "hash:old-password"}
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
	NewHandler(application.NewService(repo, fakeHasher{}, nil)).Register(group)

	req := httptest.NewRequest(http.MethodPost, "/api/system/user/password", bytes.NewBufferString(`{"oldPassword":"old-password","newPassword":"new-password"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if repo.updatedHash != "hash:new-password" {
		t.Fatalf("updated hash = %q, want hash:new-password", repo.updatedHash)
	}
	var body response.Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	data, ok := body.Data.(map[string]any)
	if !ok || data["changed"] != true {
		t.Fatalf("data = %v, want changed true", body.Data)
	}
}

func TestUpdateProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &fakeRepository{}
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
	NewHandler(application.NewService(repo, fakeHasher{}, nil)).Register(group)

	req := httptest.NewRequest(http.MethodPut, "/api/system/user/profile", bytes.NewBufferString(`{"nickName":"New Admin","headerImg":"avatar.png","phone":"10086","email":"admin@example.com","enable":2}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if repo.updatedProfile.NickName != "New Admin" || repo.updatedProfile.Email != "admin@example.com" {
		t.Fatalf("updated profile = %+v, want request fields", repo.updatedProfile)
	}
	var body response.Body
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T, want map[string]any", body.Data)
	}
	user, ok := data["user"].(map[string]any)
	if !ok || user["nickName"] != "New Admin" || user["enable"] != float64(1) {
		t.Fatalf("user = %v, want updated user without client enable override", data["user"])
	}
}

type fakeRepository struct {
	passwordHash   string
	updatedHash    string
	updatedProfile domain.ProfilePatch
}

func (r *fakeRepository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	return domain.User{}, domain.ErrUserNotFound
}

func (r *fakeRepository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.User], error) {
	return pagination.Result[domain.User]{}, domain.ErrRepositoryUnavailable
}

func (r *fakeRepository) FindPasswordHashByID(ctx context.Context, id uint) (string, error) {
	return r.passwordHash, nil
}

func (r *fakeRepository) UpdatePasswordHash(ctx context.Context, id uint, passwordHash string) error {
	r.updatedHash = passwordHash
	return nil
}

func (r *fakeRepository) UpdateProfile(ctx context.Context, id uint, profile domain.ProfilePatch) (domain.User, error) {
	r.updatedProfile = profile
	return domain.User{
		ID:          id,
		Username:    "admin",
		NickName:    profile.NickName,
		HeaderImg:   profile.HeaderImg,
		AuthorityID: 888,
		Phone:       profile.Phone,
		Email:       profile.Email,
		Enable:      1,
	}, nil
}

func (r *fakeRepository) Create(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	return domain.User{ID: 2, UUID: "new", Username: input.Username, Enable: 1}, nil
}

func (r *fakeRepository) Delete(ctx context.Context, id uint) error {
	return nil
}

func (r *fakeRepository) SetAuthorities(ctx context.Context, input domain.SetAuthoritiesInput) error {
	return nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(password string) (string, error) {
	return "hash:" + password, nil
}

func (fakeHasher) Check(password string, hash string) bool {
	return hash == "hash:"+password
}
