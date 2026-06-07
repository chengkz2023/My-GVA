package application

import (
	"context"
	"errors"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/domain"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz"
	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
)

func TestList(t *testing.T) {
	repo := &fakeRepository{apis: []domain.Api{{ID: 1, Path: "/api/test", Method: "GET", ApiGroup: "test", Description: "test api"}}}
	service := NewService(repo, nil, nil, false)

	got, err := service.List(actorContext(), domain.ListQuery{Page: pagination.Page{Page: 1, PageSize: 10}})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got.List) != 1 || got.List[0].ID != 1 {
		t.Fatalf("list = %+v, want one item", got.List)
	}
}

func TestListRepositoryUnavailable(t *testing.T) {
	got, err := NewService(&fakeRepository{err: domain.ErrRepositoryUnavailable}, nil, nil, false).List(actorContext(), domain.ListQuery{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got.List) != 0 {
		t.Fatalf("list = %+v, want empty list", got.List)
	}
}

func TestGetAll(t *testing.T) {
	repo := &fakeRepository{apis: []domain.Api{
		{ID: 1, Path: "/api/test", Method: "GET"},
		{ID: 2, Path: "/api/test", Method: "POST"},
	}}
	provider := &fakePolicyProvider{policies: []authz.Policy{{Path: "/api/test", Method: "GET"}}}
	service := NewService(repo, provider, nil, true)

	got, err := service.GetAll(actorContext())
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	if len(got.List) != 1 || got.List[0].Method != "GET" {
		t.Fatalf("list = %+v, want one filtered item", got.List)
	}
}

func TestGetAllMissingActor(t *testing.T) {
	_, err := NewService(&fakeRepository{}, nil, nil, false).GetAll(context.Background())
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Unauthorized {
		t.Fatalf("error = %v, want unauthorized app error", err)
	}
}

func TestGroups(t *testing.T) {
	repo := &fakeRepository{groups: []string{"admin", "api"}}
	got, err := NewService(repo, nil, nil, false).Groups(actorContext())
	if err != nil {
		t.Fatalf("Groups() error = %v", err)
	}
	if len(got.Groups) != 2 {
		t.Fatalf("groups = %+v, want two items", got.Groups)
	}
}

func TestPolicies(t *testing.T) {
	provider := &fakePolicyProvider{policies: []authz.Policy{{Path: "/api/test", Method: "GET"}}}
	got, err := NewService(nil, provider, nil, false).Policies(actorContext(), 888)
	if err != nil {
		t.Fatalf("Policies() error = %v", err)
	}
	if len(got.Paths) != 1 || got.Paths[0].Path != "/api/test" {
		t.Fatalf("paths = %+v, want one item", got.Paths)
	}
}

func TestUpdatePolicies(t *testing.T) {
	syncer := &fakePolicySyncer{}
	err := NewService(nil, nil, syncer, false).UpdatePolicies(actorContext(), 888, []PolicyResponse{
		{Path: "/api/test", Method: "GET"},
	})
	if err != nil {
		t.Fatalf("UpdatePolicies() error = %v", err)
	}
	if syncer.authorityID != 888 || len(syncer.policies) != 1 {
		t.Fatalf("syncer called with authorityID=%d policies=%+v", syncer.authorityID, syncer.policies)
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
	apis   []domain.Api
	groups []string
	err    error
}

func (r *fakeRepository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.Api], error) {
	if r.err != nil {
		return pagination.Result[domain.Api]{}, r.err
	}
	return pagination.Result[domain.Api]{List: r.apis, Total: int64(len(r.apis)), Page: 1, PageSize: 10}, nil
}

func (r *fakeRepository) GetAll(ctx context.Context) ([]domain.Api, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.apis, nil
}

func (r *fakeRepository) Groups(ctx context.Context) ([]string, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.groups, nil
}

type fakePolicyProvider struct {
	policies []authz.Policy
	err      error
}

func (p *fakePolicyProvider) Policies(authorityID uint) ([]authz.Policy, error) {
	return p.policies, p.err
}

type fakePolicySyncer struct {
	authorityID uint
	policies    []authz.Policy
	err         error
}

func (s *fakePolicySyncer) RemovePolicies(path, method string) error { return s.err }
func (s *fakePolicySyncer) RefreshPolicies() error { return s.err }
func (s *fakePolicySyncer) SyncPolicies(authorityID uint, policies []authz.Policy) error {
	s.authorityID = authorityID
	s.policies = policies
	return s.err
}

func (r *fakeRepository) Save(ctx context.Context, input domain.SaveApiInput) (uint, error) {
	if r.err != nil { return 0, r.err }
	return 1, nil
}
func (r *fakeRepository) Delete(ctx context.Context, id uint) (domain.SaveApiInput, error) {
	return domain.SaveApiInput{}, r.err
}
func (r *fakeRepository) DeleteByIds(ctx context.Context, ids []int) ([]domain.SaveApiInput, error) {
	return nil, r.err
}
func (r *fakeRepository) FindByID(ctx context.Context, id uint) (domain.Api, error) {
	return domain.Api{}, r.err
}
func (r *fakeRepository) GetIgnored(ctx context.Context) ([]domain.Api, error)           { return nil, r.err }
func (r *fakeRepository) Ignore(ctx context.Context, path, method string, flag bool) error { return r.err }
func (r *fakeRepository) BatchCreate(ctx context.Context, apis []domain.SaveApiInput) error { return r.err }
func (r *fakeRepository) BatchDeleteByPathMethod(ctx context.Context, apis []domain.SaveApiInput) error {
	return r.err
}
