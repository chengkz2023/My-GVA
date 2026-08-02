package mysql

import (
	"context"

	"github.com/chengkz2023/My-GVA/server/internal/modules/business/file/domain"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.File], error) {
	if r == nil || r.db == nil {
		return pagination.Result[domain.File]{}, domain.ErrRepositoryUnavailable
	}
	page := pagination.Normalize(query.Page)
	db := r.db.WithContext(ctx).Model(&ExaFileUploadAndDownload{})
	if query.ClassID > 0 {
		db = db.Where("class_id = ?", query.ClassID)
	}
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return pagination.Result[domain.File]{}, err
	}
	var entities []ExaFileUploadAndDownload
	if err := db.Order("id desc").Offset(page.Offset()).Limit(page.Limit()).Find(&entities).Error; err != nil {
		return pagination.Result[domain.File]{}, err
	}
	files := make([]domain.File, 0, len(entities))
	for _, e := range entities {
		files = append(files, mapFile(e))
	}
	return pagination.Result[domain.File]{
		List: files, Total: total, Page: page.Page, PageSize: page.PageSize,
	}, nil
}

func (r *Repository) FindByID(ctx context.Context, id uint) (domain.File, error) {
	if r == nil || r.db == nil {
		return domain.File{}, domain.ErrRepositoryUnavailable
	}
	var entity ExaFileUploadAndDownload
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		return domain.File{}, err
	}
	return mapFile(entity), nil
}

func (r *Repository) Create(ctx context.Context, file domain.File) (domain.File, error) {
	if r == nil || r.db == nil {
		return domain.File{}, domain.ErrRepositoryUnavailable
	}
	entity := ExaFileUploadAndDownload{
		Name:    file.Name,
		ClassId: file.ClassID,
		Url:     file.URL,
		Tag:     file.Tag,
		Key:     file.Key,
	}
	if err := r.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return domain.File{}, err
	}
	return mapFile(entity), nil
}

func (r *Repository) Update(ctx context.Context, id uint, name string, tag string) (domain.File, error) {
	if r == nil || r.db == nil {
		return domain.File{}, domain.ErrRepositoryUnavailable
	}
	updates := map[string]any{"name": name, "tag": tag}
	if err := r.db.WithContext(ctx).Model(&ExaFileUploadAndDownload{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return domain.File{}, err
	}
	return r.FindByID(ctx, id)
}

func (r *Repository) Delete(ctx context.Context, id uint) (domain.File, error) {
	if r == nil || r.db == nil {
		return domain.File{}, domain.ErrRepositoryUnavailable
	}
	file, err := r.FindByID(ctx, id)
	if err != nil {
		return domain.File{}, err
	}
	if err := r.db.WithContext(ctx).Delete(&ExaFileUploadAndDownload{}, id).Error; err != nil {
		return domain.File{}, err
	}
	return file, nil
}

func mapFile(e ExaFileUploadAndDownload) domain.File {
	return domain.File{
		ID:      e.ID,
		Name:    e.Name,
		ClassID: e.ClassId,
		URL:     e.Url,
		Tag:     e.Tag,
		Key:     e.Key,
	}
}
