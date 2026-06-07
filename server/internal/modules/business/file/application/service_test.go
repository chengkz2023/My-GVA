package application

import (
	"context"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/business/file/domain"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
)

func TestListFiles(t *testing.T) {
	repo := &fakeFileRepo{files: []domain.File{{ID: 1, Name: "test.png", URL: "/uploads/test.png", Tag: "png", Key: "abc.png"}}}
	service := NewService(repo, "/tmp")
	got, err := service.List(context.Background(), domain.ListQuery{Page: pagination.Page{Page: 1, PageSize: 10}})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got.Total != 1 || len(got.List) != 1 {
		t.Fatalf("result = %+v, want 1 item", got)
	}
}

func TestListFilesEmpty(t *testing.T) {
	got, err := NewService(nil, "").List(context.Background(), domain.ListQuery{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(got.List) != 0 {
		t.Fatalf("list = %+v, want empty", got.List)
	}
}

func TestUpdateFile(t *testing.T) {
	repo := &fakeFileRepo{files: []domain.File{{ID: 1, Name: "old.png", Tag: "old"}}}
	got, err := NewService(repo, "/tmp").Update(context.Background(), 1, "new.png", "png")
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if got.Name != "new.png" || got.Tag != "png" {
		t.Fatalf("file = %+v, want updated", got)
	}
}

func TestDeleteFile(t *testing.T) {
	err := NewService(&fakeFileRepo{files: []domain.File{{ID: 1, Key: "test.png"}}}, "/tmp").Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

type fakeFileRepo struct {
	files []domain.File
	err   error
}

func (r *fakeFileRepo) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.File], error) {
	if r.err != nil {
		return pagination.Result[domain.File]{}, r.err
	}
	return pagination.Result[domain.File]{List: r.files, Total: int64(len(r.files)), Page: 1, PageSize: 10}, nil
}

func (r *fakeFileRepo) FindByID(ctx context.Context, id uint) (domain.File, error) {
	if r.err != nil {
		return domain.File{}, r.err
	}
	return r.files[0], nil
}

func (r *fakeFileRepo) Create(ctx context.Context, file domain.File) (domain.File, error) {
	if r.err != nil {
		return domain.File{}, r.err
	}
	return domain.File{ID: 1, Name: file.Name, URL: file.URL, Tag: file.Tag, Key: file.Key}, nil
}

func (r *fakeFileRepo) Update(ctx context.Context, id uint, name string, tag string) (domain.File, error) {
	if r.err != nil {
		return domain.File{}, r.err
	}
	return domain.File{ID: id, Name: name, Tag: tag}, nil
}

func (r *fakeFileRepo) Delete(ctx context.Context, id uint) (domain.File, error) {
	if r.err != nil {
		return domain.File{}, r.err
	}
	return r.files[0], nil
}
