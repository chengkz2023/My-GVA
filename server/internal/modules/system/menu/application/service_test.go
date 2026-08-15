package application

import (
	"context"
	"errors"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/menu/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
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
	service := NewService(repo, nil, false)

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
	if len(got.Menus) != 1 || got.Menus[0].ID != 1 || len(got.Menus[0].Children) != 1 {
		t.Fatalf("tree = %+v, want root with one child", got.Menus)
	}
}

func TestTreeMissingActor(t *testing.T) {
	_, err := NewService(&fakeRepository{}, nil, false).Tree(context.Background())
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Unauthorized {
		t.Fatalf("error = %v, want unauthorized app error", err)
	}
}

func TestTreeRepositoryUnavailable(t *testing.T) {
	got, err := NewService(&fakeRepository{err: domain.ErrRepositoryUnavailable}, nil, false).Tree(actorContext())
	if err != nil {
		t.Fatalf("Tree() error = %v", err)
	}
	if len(got.Menus) != 0 {
		t.Fatalf("list = %+v, want empty list", got.Menus)
	}
}

func TestAssignAuthorityNonStrictSkipsScopeCheck(t *testing.T) {
	deny := AuthorityChecker(func(ctx context.Context, adminID, targetID uint) error {
		return errors.New("denied")
	})
	service := NewService(&fakeRepository{}, deny, false)

	if err := service.AssignAuthority(actorContext(), 999, []uint{1, 2}); err != nil {
		t.Fatalf("AssignAuthority() error = %v, want success in non-strict mode", err)
	}
}

func TestAssignAuthorityStrictBlocksOutOfScope(t *testing.T) {
	deny := AuthorityChecker(func(ctx context.Context, adminID, targetID uint) error {
		return errors.New("denied")
	})
	service := NewService(&fakeRepository{}, deny, true)

	err := service.AssignAuthority(actorContext(), 999, []uint{1, 2})
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Forbidden {
		t.Fatalf("error = %v, want forbidden app error", err)
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

func (r *fakeRepository) All(ctx context.Context) ([]domain.Menu, error) {
	if r.err != nil { return nil, r.err }
	return r.menus, nil
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
