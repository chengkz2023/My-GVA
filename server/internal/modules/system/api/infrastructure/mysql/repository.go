package mysql

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/api/domain"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.Api], error) {
	if r == nil || r.db == nil {
		return pagination.Result[domain.Api]{}, domain.ErrRepositoryUnavailable
	}

	db := r.db.WithContext(ctx).Model(&SysApi{})
	if query.Path != "" {
		db = db.Where("path LIKE ?", "%"+query.Path+"%")
	}
	if query.Description != "" {
		db = db.Where("description LIKE ?", "%"+query.Description+"%")
	}
	if query.Method != "" {
		db = db.Where("method = ?", query.Method)
	}
	if query.ApiGroup != "" {
		db = db.Where("api_group = ?", query.ApiGroup)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return pagination.Result[domain.Api]{}, err
	}

	page := pagination.Normalize(query.Page)
	orderKey := "id"
	if query.OrderKey != "" {
		orderKey = query.OrderKey
	}
	order := fmt.Sprintf("%s asc", orderKey)
	if query.Desc {
		order = fmt.Sprintf("%s desc", orderKey)
	}

	var entities []SysApi
	if err := db.Order(order).Offset(page.Offset()).Limit(page.Limit()).Find(&entities).Error; err != nil {
		return pagination.Result[domain.Api]{}, err
	}

	apis := make([]domain.Api, 0, len(entities))
	for _, entity := range entities {
		apis = append(apis, mapSysApi(entity))
	}
	return pagination.Result[domain.Api]{
		List:     apis,
		Total:    total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]domain.Api, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var entities []SysApi
	if err := r.db.WithContext(ctx).Order("id desc").Find(&entities).Error; err != nil {
		return nil, err
	}
	apis := make([]domain.Api, 0, len(entities))
	for _, entity := range entities {
		apis = append(apis, mapSysApi(entity))
	}
	return apis, nil
}

func (r *Repository) Groups(ctx context.Context) ([]string, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var groups []string
	if err := r.db.WithContext(ctx).Model(&SysApi{}).Distinct("api_group").Pluck("api_group", &groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *Repository) Save(ctx context.Context, input domain.SaveApiInput) (uint, error) {
	if r == nil || r.db == nil {
		return 0, domain.ErrRepositoryUnavailable
	}

	updating := input.ID != 0

	if updating {
		var old SysApi
		if err := r.db.WithContext(ctx).Where("id = ?", input.ID).First(&old).Error; err != nil {
			return 0, domain.ErrApiNotFound
		}
		if old.Path != input.Path || old.Method != input.Method {
			var count int64
			r.db.WithContext(ctx).Model(&SysApi{}).
				Where("id <> ? AND path = ? AND method = ?", input.ID, input.Path, input.Method).Count(&count)
			if count > 0 {
				return 0, domain.ErrApiDuplicate
			}
		}
	} else {
		var count int64
		r.db.WithContext(ctx).Model(&SysApi{}).
			Where("path = ? AND method = ?", input.Path, input.Method).Count(&count)
		if count > 0 {
			return 0, domain.ErrApiDuplicate
		}
	}

	api := SysApi{
		Path:        input.Path,
		Description: input.Description,
		ApiGroup:    input.ApiGroup,
		Method:      input.Method,
	}
	if updating {
		api.ID = input.ID
		if err := r.db.WithContext(ctx).Save(&api).Error; err != nil {
			return 0, err
		}
	} else {
		if err := r.db.WithContext(ctx).Create(&api).Error; err != nil {
			return 0, err
		}
	}
	return api.ID, nil
}

func (r *Repository) Delete(ctx context.Context, id uint) (domain.SaveApiInput, error) {
	if r == nil || r.db == nil {
		return domain.SaveApiInput{}, domain.ErrRepositoryUnavailable
	}
	var entity SysApi
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		return domain.SaveApiInput{}, err
	}
	if err := r.db.WithContext(ctx).Delete(&entity).Error; err != nil {
		return domain.SaveApiInput{}, err
	}
	return domain.SaveApiInput{
		ID: entity.ID, Path: entity.Path, Description: entity.Description,
		ApiGroup: entity.ApiGroup, Method: entity.Method,
	}, nil
}

func (r *Repository) DeleteByIds(ctx context.Context, ids []int) ([]domain.SaveApiInput, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var entities []SysApi
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&entities).Error; err != nil {
		return nil, err
	}
	result := make([]domain.SaveApiInput, 0, len(entities))
	for _, e := range entities {
		result = append(result, domain.SaveApiInput{
			ID: e.ID, Path: e.Path, Description: e.Description, ApiGroup: e.ApiGroup, Method: e.Method,
		})
	}
	if err := r.db.WithContext(ctx).Delete(&SysApi{}, "id IN ?", ids).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) FindByID(ctx context.Context, id uint) (domain.Api, error) {
	if r == nil || r.db == nil {
		return domain.Api{}, domain.ErrRepositoryUnavailable
	}
	var entity SysApi
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		return domain.Api{}, err
	}
	return mapSysApi(entity), nil
}

func (r *Repository) GetIgnored(ctx context.Context) ([]domain.Api, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var entities []SysIgnoreApi
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, err
	}
	apis := make([]domain.Api, 0, len(entities))
	for _, e := range entities {
		apis = append(apis, domain.Api{Path: e.Path, Method: e.Method})
	}
	return apis, nil
}

func (r *Repository) Ignore(ctx context.Context, path, method string, flag bool) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	if flag {
		return r.db.WithContext(ctx).Create(&SysIgnoreApi{Path: path, Method: method}).Error
	}
	return r.db.WithContext(ctx).Where("path = ? AND method = ?", path, method).Delete(&SysIgnoreApi{}).Error
}

func (r *Repository) BatchCreate(ctx context.Context, apis []domain.SaveApiInput) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	if len(apis) == 0 {
		return nil
	}
	entities := make([]SysApi, 0, len(apis))
	for _, a := range apis {
		entities = append(entities, SysApi{
			Path: a.Path, Description: a.Description, ApiGroup: a.ApiGroup, Method: a.Method,
		})
	}
	return r.db.WithContext(ctx).Create(&entities).Error
}

func (r *Repository) BatchDeleteByPathMethod(ctx context.Context, apis []domain.SaveApiInput) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	if len(apis) == 0 {
		return nil
	}
	for _, a := range apis {
		r.db.WithContext(ctx).Where("path = ? AND method = ?", a.Path, a.Method).Delete(&SysApi{})
	}
	return nil
}

func mapSysApi(entity SysApi) domain.Api {
	return domain.Api{
		ID:          entity.ID,
		Path:        entity.Path,
		Description: entity.Description,
		ApiGroup:    entity.ApiGroup,
		Method:      entity.Method,
	}
}
