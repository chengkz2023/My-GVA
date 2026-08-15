package mysql

import (
	"context"
	"errors"
	"time"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/user/domain"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	if r == nil || r.db == nil {
		return domain.User{}, domain.ErrRepositoryUnavailable
	}

	var user platformdb.SysUser
	err := r.db.WithContext(ctx).Select(
		"id",
		"uuid",
		"username",
		"nick_name",
		"header_img",
		"authority_id",
		"phone",
		"email",
		"enable",
	).Preload("Authorities").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, domain.ErrUserNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return mapUser(user), nil
}

func (r *Repository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.User], error) {
	page := pagination.Normalize(query.Page)
	if r == nil || r.db == nil {
		return pagination.Result[domain.User]{
			List:     []domain.User{},
			Page:     page.Page,
			PageSize: page.PageSize,
		}, domain.ErrRepositoryUnavailable
	}

	db := r.db.WithContext(ctx).Model(&platformdb.SysUser{})
	if query.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+query.NickName+"%")
	}
	if query.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+query.Phone+"%")
	}
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", "%"+query.Email+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return pagination.Result[domain.User]{}, err
	}

	var users []platformdb.SysUser
	err := db.Select(
		"id",
		"uuid",
		"username",
		"nick_name",
		"header_img",
		"authority_id",
		"phone",
		"email",
		"enable",
	).Preload("Authorities").Limit(page.Limit()).Offset(page.Offset()).Order("id desc").Find(&users).Error
	if err != nil {
		return pagination.Result[domain.User]{}, err
	}

	items := make([]domain.User, 0, len(users))
	for _, user := range users {
		items = append(items, mapUser(user))
	}
	return pagination.Result[domain.User]{
		List:     items,
		Total:    total,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, nil
}

func (r *Repository) FindPasswordHashByID(ctx context.Context, id uint) (string, error) {
	if r == nil || r.db == nil {
		return "", domain.ErrRepositoryUnavailable
	}

	var user platformdb.SysUser
	err := r.db.WithContext(ctx).Select("id", "password").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", domain.ErrUserNotFound
	}
	if err != nil {
		return "", err
	}
	return user.Password, nil
}

func (r *Repository) UpdatePasswordHash(ctx context.Context, id uint, passwordHash string) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	return r.db.WithContext(ctx).Model(&platformdb.SysUser{}).Where("id = ?", id).Update("password", passwordHash).Error
}

func (r *Repository) UpdateProfile(ctx context.Context, id uint, profile domain.ProfilePatch) (domain.User, error) {
	if r == nil || r.db == nil {
		return domain.User{}, domain.ErrRepositoryUnavailable
	}

	result := r.db.WithContext(ctx).Model(&platformdb.SysUser{}).Where("id = ?", id).Updates(map[string]any{
		"updated_at": time.Now(),
		"nick_name":  profile.NickName,
		"header_img": profile.HeaderImg,
		"phone":      profile.Phone,
		"email":      profile.Email,
	})
	if result.Error != nil {
		return domain.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return domain.User{}, domain.ErrUserNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *Repository) Create(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	if r == nil || r.db == nil {
		return domain.User{}, domain.ErrRepositoryUnavailable
	}

	authorities := make([]platformdb.SysAuthority, 0, len(input.AuthorityIDs))
	for _, id := range input.AuthorityIDs {
		authorities = append(authorities, platformdb.SysAuthority{AuthorityId: id})
	}

	user := platformdb.SysUser{
		UUID:        uuid.New(),
		Username:    input.Username,
		Password:    input.PasswordHash,
		NickName:    input.NickName,
		HeaderImg:   input.HeaderImg,
		AuthorityId: input.AuthorityID,
		Authorities: authorities,
		Enable:      input.Enable,
		Phone:       input.Phone,
		Email:       input.Email,
	}

	if input.NickName == "" {
		user.NickName = input.Username
	}
	if input.HeaderImg == "" {
		user.HeaderImg = "https://qmplusimg.henrongyi.top/gva_header.jpg"
	}
	if input.Enable == 0 {
		user.Enable = 1
	}

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return domain.User{}, err
	}
	return r.FindByID(ctx, user.ID)
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&platformdb.SysUser{}).Error; err != nil {
			return err
		}
		return tx.Where("sys_user_id = ?", id).Delete(&platformdb.SysUserAuthority{}).Error
	})
}

func (r *Repository) SetAuthorities(ctx context.Context, input domain.SetAuthoritiesInput) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user platformdb.SysUser
		if err := tx.First(&user, input.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrUserNotFound
			}
			return err
		}

		if err := tx.Where("sys_user_id = ?", input.UserID).Delete(&platformdb.SysUserAuthority{}).Error; err != nil {
			return err
		}

		if len(input.AuthorityIDs) > 0 {
			userAuthorities := make([]platformdb.SysUserAuthority, 0, len(input.AuthorityIDs))
			for _, authID := range input.AuthorityIDs {
				userAuthorities = append(userAuthorities, platformdb.SysUserAuthority{
					SysUserId:               input.UserID,
					SysAuthorityAuthorityId: authID,
				})
			}
			if err := tx.Create(&userAuthorities).Error; err != nil {
				return err
			}

			if err := tx.Model(&user).Update("authority_id", input.AuthorityIDs[0]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func mapUser(user platformdb.SysUser) domain.User {
	authorityIDs := make([]uint, 0, len(user.Authorities))
	for _, a := range user.Authorities {
		authorityIDs = append(authorityIDs, a.AuthorityId)
	}
	return domain.User{
		ID:           user.ID,
		UUID:         user.UUID.String(),
		Username:     user.Username,
		NickName:     user.NickName,
		HeaderImg:    user.HeaderImg,
		AuthorityID:  user.AuthorityId,
		AuthorityIDs: authorityIDs,
		Phone:        user.Phone,
		Email:        user.Email,
		Enable:       user.Enable,
	}
}
