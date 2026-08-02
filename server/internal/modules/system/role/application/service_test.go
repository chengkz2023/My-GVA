package application

import (
	"context"
	"errors"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/role/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
)

func TestTree(t *testing.T) {
	repo := &fakeRepository{roles: []domain.Role{{
		AuthorityID:   888,
		AuthorityName: "admin",
		DefaultRouter: "dashboard",
		Children: []domain.Role{{
			AuthorityID:   889,
			AuthorityName: "sub-admin",
		}},
	}}}
	service := NewService(repo, nil, nil, true)

	got, err := service.Tree(actorContext())
	if err != nil {
		t.Fatalf("Tree() error = %v", err)
	}
	if !repo.called {
		t.Fatal("repository was not called")
	}
	if repo.authorityID != 888 || !repo.strict {
		t.Fatalf("query = %d/%v, want 888/true", repo.authorityID, repo.strict)
	}
	if len(got.List) != 1 || got.List[0].AuthorityID != 888 || len(got.List[0].Children) != 1 {
		t.Fatalf("tree = %+v, want root with one child", got.List)
	}
}

func TestTreeMissingActor(t *testing.T) {
	_, err := NewService(&fakeRepository{}, nil, nil, false).Tree(context.Background())
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Unauthorized {
		t.Fatalf("error = %v, want unauthorized app error", err)
	}
}

func TestTreeRepositoryUnavailable(t *testing.T) {
	got, err := NewService(&fakeRepository{err: domain.ErrRepositoryUnavailable}, nil, nil, false).Tree(actorContext())
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
	roles       []domain.Role
	err         error
	called      bool
	authorityID uint
	strict      bool
}

func (r *fakeRepository) Tree(ctx context.Context, authorityID uint, strict bool) ([]domain.Role, error) {
	r.called = true
	r.authorityID = authorityID
	r.strict = strict
	if r.err != nil {
		return nil, r.err
	}
	return r.roles, nil
}

func (r *fakeRepository) Save(ctx context.Context, input domain.SaveRoleInput) error { return r.err }
func (r *fakeRepository) Delete(ctx context.Context, id uint) (domain.SaveRoleInput, error) {
	return domain.SaveRoleInput{}, r.err
}
func (r *fakeRepository) FindByID(ctx context.Context, id uint) (domain.Role, error) {
	return domain.Role{AuthorityID: id, AuthorityName: "test"}, r.err
}
func (r *fakeRepository) FindMenuIDs(ctx context.Context, id uint) ([]uint, error) {
	return nil, r.err
}
func (r *fakeRepository) CopyMenusAndButtons(ctx context.Context, oldID, newID uint) error {
	return r.err
}
func (r *fakeRepository) SetDataAuthority(ctx context.Context, input domain.DataAuthorityInput) error {
	return r.err
}
func (r *fakeRepository) GetDataAuthorities(ctx context.Context, id uint) ([]uint, error) {
	return nil, r.err
}

func (r *fakeRepository) GetDescendantIDs(ctx context.Context, authorityID uint) ([]uint, error) {
	return nil, r.err
}
func (r *fakeRepository) CheckAuthorityAuth(ctx context.Context, adminID, targetID uint) error {
	return r.err
}
