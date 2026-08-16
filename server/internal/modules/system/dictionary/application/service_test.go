package application

import (
	"context"
	"errors"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

func TestCreateRequiresTypeAndName(t *testing.T) {
	_, err := NewService(&fakeRepository{}).Create(actorContext(), SaveDictionaryCommand{Type: "", Name: "x"})
	assertKind(t, err, apperrors.Validation)
}

func TestCreateDuplicateType(t *testing.T) {
	_, err := NewService(&fakeRepository{err: domain.ErrTypeExists}).Create(actorContext(), SaveDictionaryCommand{Type: "gender", Name: "性别"})
	assertKind(t, err, apperrors.Conflict)
}

func TestCreateSuccess(t *testing.T) {
	got, err := NewService(&fakeRepository{id: 3}).Create(actorContext(), SaveDictionaryCommand{Type: "gender", Name: "性别", Sort: 1})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if got.ID != 3 || got.Type != "gender" || got.Status != 1 {
		t.Fatalf("got = %+v, want id 3 gender status 1", got)
	}
}

func TestListRequiresActor(t *testing.T) {
	_, err := NewService(&fakeRepository{}).List(context.Background(), ListDictionaryQuery{})
	assertKind(t, err, apperrors.Unauthorized)
}

func TestListRepositoryUnavailable(t *testing.T) {
	got, err := NewService(&fakeRepository{err: domain.ErrRepositoryUnavailable}).List(actorContext(), ListDictionaryQuery{Page: pagination.Page{Page: 1, PageSize: 10}})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got.List) != 0 || got.Total != 0 {
		t.Fatalf("got = %+v, want empty list", got)
	}
}

func TestDeleteNotFound(t *testing.T) {
	err := NewService(&fakeRepository{err: domain.ErrDictionaryNotFound}).Delete(actorContext(), 1)
	assertKind(t, err, apperrors.NotFound)
}

func TestCreateDetailRequiresFields(t *testing.T) {
	_, err := NewService(&fakeRepository{}).CreateDetail(actorContext(), SaveDetailCommand{DictionaryID: 1, Label: "", Value: ""})
	assertKind(t, err, apperrors.Validation)
}

func TestCreateDetailDictionaryMissing(t *testing.T) {
	_, err := NewService(&fakeRepository{err: domain.ErrDictionaryNotFound}).CreateDetail(actorContext(), SaveDetailCommand{DictionaryID: 9, Label: "男", Value: "male"})
	assertKind(t, err, apperrors.NotFound)
}

func TestTypesEmptyWhenRepoUnavailable(t *testing.T) {
	got, err := NewService(&fakeRepository{err: domain.ErrRepositoryUnavailable}).Types(actorContext())
	if err != nil {
		t.Fatalf("Types() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got = %+v, want empty", got)
	}
}

func TestTypesReturnsDictionaryWithDetails(t *testing.T) {
	repo := &fakeRepository{
		types: []domain.DictionaryWithDetails{{
			Dictionary: domain.Dictionary{ID: 1, Type: "gender", Name: "性别", Status: 1},
			Details: []domain.DictionaryDetail{
				{ID: 1, DictionaryID: 1, Label: "男", Value: "male", Status: 1},
			},
		}},
	}
	got, err := NewService(repo).Types(actorContext())
	if err != nil {
		t.Fatalf("Types() error = %v", err)
	}
	if len(got) != 1 || got[0].Type != "gender" || len(got[0].Details) != 1 || got[0].Details[0].Value != "male" {
		t.Fatalf("got = %+v, want gender dictionary with one detail", got)
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

func assertKind(t *testing.T, err error, kind apperrors.Kind) {
	t.Helper()
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != kind {
		t.Fatalf("error = %v, want kind %v", err, kind)
	}
}

type fakeRepository struct {
	id    uint
	err   error
	types []domain.DictionaryWithDetails
}

func (r *fakeRepository) List(ctx context.Context, query domain.ListDictionaryQuery) (pagination.Result[domain.Dictionary], error) {
	if r.err != nil {
		return pagination.Result[domain.Dictionary]{}, r.err
	}
	return pagination.Result[domain.Dictionary]{List: []domain.Dictionary{}, Page: query.Page.Page, PageSize: query.Page.PageSize}, nil
}

func (r *fakeRepository) FindByID(ctx context.Context, id uint) (domain.Dictionary, error) {
	if r.err != nil {
		return domain.Dictionary{}, r.err
	}
	return domain.Dictionary{ID: id}, nil
}

func (r *fakeRepository) Save(ctx context.Context, input domain.SaveDictionaryInput) (uint, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.id != 0 {
		return r.id, nil
	}
	return 1, nil
}

func (r *fakeRepository) Delete(ctx context.Context, id uint) error { return r.err }

func (r *fakeRepository) ListDetails(ctx context.Context, dictionaryID uint) ([]domain.DictionaryDetail, error) {
	if r.err != nil {
		return nil, r.err
	}
	return []domain.DictionaryDetail{}, nil
}

func (r *fakeRepository) SaveDetail(ctx context.Context, input domain.SaveDetailInput) (uint, error) {
	if r.err != nil {
		return 0, r.err
	}
	return 1, nil
}

func (r *fakeRepository) DeleteDetail(ctx context.Context, id uint) error { return r.err }

func (r *fakeRepository) Types(ctx context.Context) ([]domain.DictionaryWithDetails, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.types, nil
}
