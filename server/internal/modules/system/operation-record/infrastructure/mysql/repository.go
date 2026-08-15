package mysql

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/operation-record/domain"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.Record], error) {
	if r == nil || r.db == nil {
		return pagination.Result[domain.Record]{}, domain.ErrRepositoryUnavailable
	}
	page := pagination.Normalize(query.Page)
	db := r.db.WithContext(ctx).Model(&SysOperationRecord{}).Preload("User")
	if query.Method != "" {
		db = db.Where("method = ?", query.Method)
	}
	if query.Path != "" {
		db = db.Where("path LIKE ?", "%"+query.Path+"%")
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return pagination.Result[domain.Record]{}, err
	}
	var records []SysOperationRecord
	if err := db.Order("id desc").Offset(page.Offset()).Limit(page.Limit()).Find(&records).Error; err != nil {
		return pagination.Result[domain.Record]{}, err
	}
	items := make([]domain.Record, 0, len(records))
	for _, r := range records {
		items = append(items, mapRecord(r))
	}
	return pagination.Result[domain.Record]{
		List: items, Total: total, Page: page.Page, PageSize: page.PageSize,
	}, nil
}

func (r *Repository) FindByID(ctx context.Context, id uint) (domain.Record, error) {
	if r == nil || r.db == nil {
		return domain.Record{}, domain.ErrRepositoryUnavailable
	}
	var record SysOperationRecord
	if err := r.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Record{}, domain.ErrRecordNotFound
		}
		return domain.Record{}, err
	}
	return mapRecord(record), nil
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	return r.db.WithContext(ctx).Delete(&SysOperationRecord{}, id).Error
}

func (r *Repository) DeleteByIds(ctx context.Context, ids []int) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	return r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&SysOperationRecord{}).Error
}

func mapRecord(r SysOperationRecord) domain.Record {
	rec := domain.Record{
		ID: r.ID, CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
		IP: r.Ip, Method: r.Method, Path: r.Path,
		Status: r.Status, Latency: int64(r.Latency), Agent: r.Agent,
		ErrorMessage: r.ErrorMessage, Body: r.Body, Resp: r.Resp,
		UserID: r.UserID,
	}
	rec.Username = r.User.Username
	rec.NickName = r.User.NickName
	return rec
}
