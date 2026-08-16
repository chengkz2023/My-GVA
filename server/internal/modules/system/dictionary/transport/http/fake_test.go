package http

import (
	"context"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/domain"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

type fakeRepository struct{}

func (fakeRepository) List(ctx context.Context, query domain.ListDictionaryQuery) (pagination.Result[domain.Dictionary], error) {
	return pagination.Result[domain.Dictionary]{List: []domain.Dictionary{}, Page: 1, PageSize: 10}, nil
}
func (fakeRepository) FindByID(ctx context.Context, id uint) (domain.Dictionary, error) {
	return domain.Dictionary{}, nil
}
func (fakeRepository) Save(ctx context.Context, input domain.SaveDictionaryInput) (uint, error) {
	return 1, nil
}
func (fakeRepository) Delete(ctx context.Context, id uint) error { return nil }
func (fakeRepository) ListDetails(ctx context.Context, dictionaryID uint) ([]domain.DictionaryDetail, error) {
	return []domain.DictionaryDetail{}, nil
}
func (fakeRepository) SaveDetail(ctx context.Context, input domain.SaveDetailInput) (uint, error) {
	return 1, nil
}
func (fakeRepository) DeleteDetail(ctx context.Context, id uint) error { return nil }
func (fakeRepository) Types(ctx context.Context) ([]domain.DictionaryWithDetails, error) {
	return []domain.DictionaryWithDetails{}, nil
}
