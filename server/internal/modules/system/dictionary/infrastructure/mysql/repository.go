package mysql

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/dictionary/domain"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, query domain.ListDictionaryQuery) (pagination.Result[domain.Dictionary], error) {
	if r == nil || r.db == nil {
		return pagination.Result[domain.Dictionary]{}, domain.ErrRepositoryUnavailable
	}
	db := r.db.WithContext(ctx).Model(&SysDictionary{})
	if query.Type != "" {
		db = db.Where("type LIKE ?", "%"+query.Type+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return pagination.Result[domain.Dictionary]{}, err
	}

	page := pagination.Normalize(query.Page)
	var entities []SysDictionary
	if err := db.Order("sort asc").Order("id desc").
		Offset(page.Offset()).Limit(page.Limit()).Find(&entities).Error; err != nil {
		return pagination.Result[domain.Dictionary]{}, err
	}

	items := make([]domain.Dictionary, 0, len(entities))
	for _, e := range entities {
		items = append(items, mapDictionary(e))
	}
	return pagination.Result[domain.Dictionary]{
		List:     items,
		Total:    total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, nil
}

func (r *Repository) FindByID(ctx context.Context, id uint) (domain.Dictionary, error) {
	if r == nil || r.db == nil {
		return domain.Dictionary{}, domain.ErrRepositoryUnavailable
	}
	var entity SysDictionary
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Dictionary{}, domain.ErrDictionaryNotFound
		}
		return domain.Dictionary{}, err
	}
	return mapDictionary(entity), nil
}

func (r *Repository) Save(ctx context.Context, input domain.SaveDictionaryInput) (uint, error) {
	if r == nil || r.db == nil {
		return 0, domain.ErrRepositoryUnavailable
	}

	if input.ID != 0 {
		var entity SysDictionary
		if err := r.db.WithContext(ctx).First(&entity, input.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, domain.ErrDictionaryNotFound
			}
			return 0, err
		}
		if entity.Type != input.Type {
			var count int64
			if err := r.db.WithContext(ctx).Model(&SysDictionary{}).
				Where("type = ? AND id <> ?", input.Type, input.ID).Count(&count).Error; err != nil {
				return 0, err
			}
			if count > 0 {
				return 0, domain.ErrTypeExists
			}
		}
		entity.Type = input.Type
		entity.Name = input.Name
		entity.Sort = input.Sort
		entity.Status = input.Status
		if err := r.db.WithContext(ctx).Save(&entity).Error; err != nil {
			return 0, err
		}
		return entity.ID, nil
	}

	entity := SysDictionary{
		Type:   input.Type,
		Name:   input.Name,
		Sort:   input.Sort,
		Status: input.Status,
	}
	if err := r.db.WithContext(ctx).Create(&entity).Error; err != nil {
		if isDuplicateKey(err) {
			return 0, domain.ErrTypeExists
		}
		return 0, err
	}
	return entity.ID, nil
}

// Delete 删除字典及其全部字典项（事务）。
func (r *Repository) Delete(ctx context.Context, id uint) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var entity SysDictionary
		if err := tx.First(&entity, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrDictionaryNotFound
			}
			return err
		}
		if err := tx.Where("dictionary_id = ?", id).Delete(&SysDictionaryDetail{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entity).Error
	})
}

func (r *Repository) ListDetails(ctx context.Context, dictionaryID uint) ([]domain.DictionaryDetail, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var entities []SysDictionaryDetail
	if err := r.db.WithContext(ctx).
		Where("dictionary_id = ?", dictionaryID).
		Order("sort asc").Order("id asc").
		Find(&entities).Error; err != nil {
		return nil, err
	}
	items := make([]domain.DictionaryDetail, 0, len(entities))
	for _, e := range entities {
		items = append(items, mapDetail(e))
	}
	return items, nil
}

func (r *Repository) SaveDetail(ctx context.Context, input domain.SaveDetailInput) (uint, error) {
	if r == nil || r.db == nil {
		return 0, domain.ErrRepositoryUnavailable
	}
	if input.DictionaryID == 0 {
		return 0, domain.ErrDictionaryNotFound
	}

	if input.ID != 0 {
		var entity SysDictionaryDetail
		if err := r.db.WithContext(ctx).First(&entity, input.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, domain.ErrDictionaryDetailNotFound
			}
			return 0, err
		}
		entity.Label = input.Label
		entity.Value = input.Value
		entity.Sort = input.Sort
		entity.Status = input.Status
		if err := r.db.WithContext(ctx).Save(&entity).Error; err != nil {
			return 0, err
		}
		return entity.ID, nil
	}

	// 校验所属字典存在
	var count int64
	if err := r.db.WithContext(ctx).Model(&SysDictionary{}).Where("id = ?", input.DictionaryID).Count(&count).Error; err != nil {
		return 0, err
	}
	if count == 0 {
		return 0, domain.ErrDictionaryNotFound
	}

	entity := SysDictionaryDetail{
		DictionaryID: input.DictionaryID,
		Label:        input.Label,
		Value:        input.Value,
		Sort:         input.Sort,
		Status:       input.Status,
	}
	if err := r.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (r *Repository) DeleteDetail(ctx context.Context, id uint) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	res := r.db.WithContext(ctx).Delete(&SysDictionaryDetail{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrDictionaryDetailNotFound
	}
	return nil
}

func (r *Repository) Types(ctx context.Context) ([]domain.DictionaryWithDetails, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var dictionaries []SysDictionary
	if err := r.db.WithContext(ctx).Where("status = ?", 1).Order("sort asc").Find(&dictionaries).Error; err != nil {
		return nil, err
	}
	var details []SysDictionaryDetail
	if err := r.db.WithContext(ctx).Where("status = ?", 1).Order("sort asc").Order("id asc").Find(&details).Error; err != nil {
		return nil, err
	}
	byID := make(map[uint][]domain.DictionaryDetail, len(details))
	for _, d := range details {
		byID[d.DictionaryID] = append(byID[d.DictionaryID], mapDetail(d))
	}
	result := make([]domain.DictionaryWithDetails, 0, len(dictionaries))
	for _, d := range dictionaries {
		items := byID[d.ID]
		if items == nil {
			items = []domain.DictionaryDetail{}
		}
		result = append(result, domain.DictionaryWithDetails{
			Dictionary: mapDictionary(d),
			Details:    items,
		})
	}
	return result, nil
}

func mapDictionary(e SysDictionary) domain.Dictionary {
	return domain.Dictionary{
		ID:     e.ID,
		Type:   e.Type,
		Name:   e.Name,
		Sort:   e.Sort,
		Status: e.Status,
	}
}

func mapDetail(e SysDictionaryDetail) domain.DictionaryDetail {
	return domain.DictionaryDetail{
		ID:           e.ID,
		DictionaryID: e.DictionaryID,
		Label:        e.Label,
		Value:        e.Value,
		Sort:         e.Sort,
		Status:       e.Status,
	}
}

func isDuplicateKey(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
