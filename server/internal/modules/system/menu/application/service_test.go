package application

import (
	"context"
	"errors"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/menu/domain"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
)

func TestTree(t *testing.T) {
	repo := &fakeRepository{menus: []domain.Menu{{
		ID:        1,
		ParentID:  0,
		Path:      "/dashboard",
		Name:      "superAdmin",
		Component: "dashboard/index",
		Sort:      1,
		Title:     "Dashboard",
		Icon:      "home",
		Children: []domain.Menu{{
			ID:        2,
			ParentID:  1,
			Path:      "authority",
			Name:      "authority",
			Component: "authority/index",
			Sort:      1,
			Title:     "角色管理",
		}},
	}}}
	service := NewService(repo, nil)

	got, err := service.Tree(actorContext())
	if err != nil {
		t.Fatalf("Tree() error = %v", err)
	}
	if !repo.called {
		t.Fatal("repository was not called")
	}
	if repo.authorityID != 888 {
		t.Fatalf("authorityID = %d, want 888", repo.authorityID)
	}
	if len(got.List) != 1 || got.List[0].ID != 1 || len(got.List[0].Children) != 1 {
		t.Fatalf("tree = %+v, want root with one child", got.List)
	}
}

func TestTreeMissingActor(t *testing.T) {
	_, err := NewService(&fakeRepository{}, nil).Tree(context.Background())
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Unauthorized {
		t.Fatalf("error = %v, want unauthorized app error", err)
	}
}

func TestTreeRepositoryUnavailable(t *testing.T) {
	got, err := NewService(&fakeRepository{err: domain.ErrRepositoryUnavailable}, nil).Tree(actorContext())
	if err != nil {
		t.Fatalf("Tree() error = %v", err)
	}
	if len(got.List) != 0 {
		t.Fatalf("list = %+v, want empty list", got.List)
	}
}

func actorContext() context.Context {
	return platformauth.ContextWithActor(context.Background(), platformauth.Actor{
		UserID:      1,
		AuthorityID: 888,
		Username:    "admin",
		NickName:    "Admin",
	})
}

type fakeRepository struct {
	menus       []domain.Menu
	err         error
	called      bool
	authorityID uint
}

func (r *fakeRepository) TreeByAuthority(ctx context.Context, authorityID uint) ([]domain.Menu, error) {
	r.called = true
	r.authorityID = authorityID
	if r.err != nil {
		return nil, r.err
	}
	return r.menus, nil
}

func (r *fakeRepository) Save(ctx context.Context, input domain.SaveMenuInput) (uint, error) { return 0, r.err }
func (r *fakeRepository) Delete(ctx context.Context, id uint) error { return r.err }
func (r *fakeRepository) FindByID(ctx context.Context, id uint) (domain.MenuDetail, error) { return domain.MenuDetail{}, r.err }
func (r *fakeRepository) AssignMenus(ctx context.Context, authorityID uint, menuIDs []uint) error {
	return r.err
}
