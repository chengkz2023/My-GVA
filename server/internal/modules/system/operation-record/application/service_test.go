package application

import (
	"context"
	"errors"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/operation-record/domain"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

func TestListRecords(t *testing.T) {
	repo := &fakeRepo{records: []domain.Record{{ID: 1, IP: "127.0.0.1", Method: "GET", Path: "/api/test"}}}
	service := NewService(repo)
	got, err := service.List(context.Background(), domain.ListQuery{Page: pagination.Page{Page: 1, PageSize: 10}})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got.Total != 1 || len(got.List) != 1 {
		t.Fatalf("result = %+v, want 1 item", got)
	}
}

func TestListRecordsUnavailable(t *testing.T) {
	got, err := NewService(&fakeRepo{err: domain.ErrRepositoryUnavailable}).List(context.Background(), domain.ListQuery{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got.List) != 0 {
		t.Fatalf("list = %+v, want empty", got.List)
	}
}

func TestFindByID(t *testing.T) {
	repo := &fakeRepo{records: []domain.Record{{ID: 1, IP: "1.2.3.4", Path: "/api/test"}}}
	got, err := NewService(repo).FindByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if got.ID != 1 || got.IP != "1.2.3.4" {
		t.Fatalf("record = %+v, want id=1", got)
	}
}

func TestDeleteRecord(t *testing.T) {
	err := NewService(&fakeRepo{}).Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestDeleteByIds(t *testing.T) {
	err := NewService(&fakeRepo{}).DeleteByIds(context.Background(), []int{1, 2, 3})
	if err != nil {
		t.Fatalf("DeleteByIds() error = %v", err)
	}
}

type fakeRepo struct {
	records []domain.Record
	err     error
}

func (r *fakeRepo) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.Record], error) {
	if r.err != nil {
		return pagination.Result[domain.Record]{}, r.err
	}
	return pagination.Result[domain.Record]{List: r.records, Total: int64(len(r.records)), Page: 1, PageSize: 10}, nil
}

func (r *fakeRepo) FindByID(ctx context.Context, id uint) (domain.Record, error) {
	if r.err != nil {
		return domain.Record{}, r.err
	}
	if len(r.records) > 0 {
		return r.records[0], nil
	}
	return domain.Record{}, errors.New("not found")
}

func (r *fakeRepo) Delete(ctx context.Context, id uint) error {
	return r.err
}

func (r *fakeRepo) DeleteByIds(ctx context.Context, ids []int) error {
	return r.err
}
