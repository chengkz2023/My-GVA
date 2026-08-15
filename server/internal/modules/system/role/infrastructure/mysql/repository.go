package mysql

import (
	"context"
	"errors"
	"strconv"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/role/domain"
	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Tree(ctx context.Context, authorityID uint, strict bool) ([]domain.Role, error) {
	if r == nil || r.db == nil {
		return []domain.Role{}, domain.ErrRepositoryUnavailable
	}

	query := r.db.WithContext(ctx).Model(&platformdb.SysAuthority{})
	if strict {
		current, err := r.find(ctx, authorityID)
		if err != nil {
			return nil, err
		}
		if isRoot(current.ParentId) {
			query = query.Where("authority_id = ?", authorityID)
		} else {
			query = query.Where("parent_id = ?", authorityID)
		}
	} else {
		query = query.Where("parent_id = ?", 0)
	}

	var authorities []platformdb.SysAuthority
	if err := query.Order("authority_id asc").Find(&authorities).Error; err != nil {
		return nil, err
	}

	roles := make([]domain.Role, 0, len(authorities))
	for _, authority := range authorities {
		role, err := r.mapWithChildren(ctx, authority)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *Repository) Save(ctx context.Context, input domain.SaveRoleInput) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}

	updating := false
	var existing platformdb.SysAuthority
	err := r.db.WithContext(ctx).Where("authority_id = ?", input.AuthorityID).First(&existing).Error
	if err == nil {
		updating = true
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	if !updating {
		if input.ParentID == nil || *input.ParentID == 0 {
		}
		auth := platformdb.SysAuthority{
			AuthorityId:   input.AuthorityID,
			AuthorityName: input.AuthorityName,
			ParentId:      input.ParentID,
			DefaultRouter: input.DefaultRouter,
		}
		if err := r.db.WithContext(ctx).Create(&auth).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return domain.ErrRoleIDExists
			}
			return err
		}
		// Assign default dashboard menu
		r.assignMenu(ctx, input.AuthorityID, 1)
		return nil
	}

	updates := map[string]any{
		"authority_name": input.AuthorityName,
		"parent_id":      input.ParentID,
		"default_router": input.DefaultRouter,
	}
	if err := r.db.WithContext(ctx).Model(&existing).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, authorityID uint) (domain.SaveRoleInput, error) {
	if r == nil || r.db == nil {
		return domain.SaveRoleInput{}, domain.ErrRepositoryUnavailable
	}

	var auth platformdb.SysAuthority
	err := r.db.WithContext(ctx).
		Preload("Users").
		Where("authority_id = ?", authorityID).
		First(&auth).Error
	if err != nil {
		return domain.SaveRoleInput{}, domain.ErrRoleNotFound
	}

	if len(auth.Users) > 0 {
		return domain.SaveRoleInput{}, domain.ErrRoleHasUsers
	}

	var childCount int64
	r.db.WithContext(ctx).Model(&platformdb.SysAuthority{}).Where("parent_id = ?", authorityID).Count(&childCount)
	if childCount > 0 {
		return domain.SaveRoleInput{}, domain.ErrRoleHasChildren
	}

	result := domain.SaveRoleInput{
		AuthorityID:   auth.AuthorityId,
		AuthorityName: auth.AuthorityName,
		ParentID:      auth.ParentId,
		DefaultRouter: auth.DefaultRouter,
	}

	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var del platformdb.SysAuthority
		if err := tx.Where("authority_id = ?", authorityID).First(&del).Error; err != nil {
			return err
		}
		// Remove menu associations
		tx.Where("sys_authority_authority_id = ?", authorityID).Delete(&platformdb.SysAuthorityMenu{})
		// Remove data authority associations
		tx.Exec("DELETE FROM sys_data_authority_id WHERE sys_authority_id = ?", authorityID)
		// Remove user associations
		tx.Where("sys_authority_authority_id = ?", authorityID).Delete(&platformdb.SysUserAuthority{})
		// Remove authority buttons
		tx.Where("authority_id = ?", authorityID).Delete(&platformdb.SysAuthorityBtn{})
		// Hard delete the authority
		tx.Unscoped().Delete(&del)
		return nil
	})
	return result, err
}

func (r *Repository) FindByID(ctx context.Context, authorityID uint) (domain.Role, error) {
	if r == nil || r.db == nil {
		return domain.Role{}, domain.ErrRepositoryUnavailable
	}
	auth, err := r.find(ctx, authorityID)
	if err != nil {
		return domain.Role{}, err
	}
	return mapRole(auth), nil
}

func (r *Repository) FindMenuIDs(ctx context.Context, authorityID uint) ([]uint, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var menus []platformdb.SysAuthorityMenu
	if err := r.db.WithContext(ctx).Where("sys_authority_authority_id = ?", authorityID).Find(&menus).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(menus))
	for _, m := range menus {
		id, err := strconv.ParseUint(m.MenuId, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}
	return ids, nil
}

func (r *Repository) CopyMenusAndButtons(ctx context.Context, oldAuthorityID, newAuthorityID uint) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Copy menu assignments
		var oldMenus []platformdb.SysAuthorityMenu
		if err := tx.Where("sys_authority_authority_id = ?", oldAuthorityID).Find(&oldMenus).Error; err != nil {
			return err
		}
		for _, m := range oldMenus {
			tx.Create(&platformdb.SysAuthorityMenu{
				MenuId:      m.MenuId,
				AuthorityId: strconv.FormatUint(uint64(newAuthorityID), 10),
			})
		}
		// Copy authority buttons
		var oldBtns []platformdb.SysAuthorityBtn
		if err := tx.Where("authority_id = ?", oldAuthorityID).Find(&oldBtns).Error; err != nil {
			return err
		}
		for _, b := range oldBtns {
			tx.Create(&platformdb.SysAuthorityBtn{
				AuthorityId:      newAuthorityID,
				SysMenuID:        b.SysMenuID,
				SysBaseMenuBtnID: b.SysBaseMenuBtnID,
			})
		}
		return nil
	})
}

func (r *Repository) SetDataAuthority(ctx context.Context, input domain.DataAuthorityInput) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}
	var auth platformdb.SysAuthority
	if err := r.db.WithContext(ctx).Where("authority_id = ?", input.AuthorityID).First(&auth).Error; err != nil {
		return err
	}
	dataAuths := make([]platformdb.SysAuthority, 0, len(input.DataAuthorityIDs))
	for _, id := range input.DataAuthorityIDs {
		dataAuths = append(dataAuths, platformdb.SysAuthority{AuthorityId: id})
	}
	return r.db.WithContext(ctx).Model(&auth).Association("DataAuthorityId").Replace(dataAuths)
}

func (r *Repository) GetDataAuthorities(ctx context.Context, authorityID uint) ([]uint, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var auth platformdb.SysAuthority
	if err := r.db.WithContext(ctx).Preload("DataAuthorityId").Where("authority_id = ?", authorityID).First(&auth).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(auth.DataAuthorityId))
	for _, da := range auth.DataAuthorityId {
		ids = append(ids, da.AuthorityId)
	}
	return ids, nil
}

func (r *Repository) find(ctx context.Context, authorityID uint) (platformdb.SysAuthority, error) {
	var authority platformdb.SysAuthority
	err := r.db.WithContext(ctx).Where("authority_id = ?", authorityID).First(&authority).Error
	return authority, err
}

func (r *Repository) mapWithChildren(ctx context.Context, authority platformdb.SysAuthority) (domain.Role, error) {
	role := mapRole(authority)

	var children []platformdb.SysAuthority
	err := r.db.WithContext(ctx).Where("parent_id = ?", authority.AuthorityId).Order("authority_id asc").Find(&children).Error
	if err != nil {
		return domain.Role{}, err
	}
	role.Children = make([]domain.Role, 0, len(children))
	for _, child := range children {
		childRole, err := r.mapWithChildren(ctx, child)
		if err != nil {
			return domain.Role{}, err
		}
		role.Children = append(role.Children, childRole)
	}
	return role, nil
}

func (r *Repository) assignMenu(ctx context.Context, authorityID, menuID uint) {
	var auth platformdb.SysAuthority
	auth.AuthorityId = authorityID
	menus := []platformdb.SysBaseMenu{{}}
	menus[0].ID = menuID
	_ = r.db.WithContext(ctx).Model(&auth).Association("SysBaseMenus").Append(menus)
}

func mapRole(authority platformdb.SysAuthority) domain.Role {
	return domain.Role{
		AuthorityID:   authority.AuthorityId,
		AuthorityName: authority.AuthorityName,
		ParentID:      authority.ParentId,
		DefaultRouter: authority.DefaultRouter,
	}
}

func (r *Repository) GetDescendantIDs(ctx context.Context, authorityID uint) ([]uint, error) {
	if r == nil || r.db == nil {
		return nil, domain.ErrRepositoryUnavailable
	}
	var ids []uint
	if err := r.collectDescendants(ctx, authorityID, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *Repository) collectDescendants(ctx context.Context, parentID uint, ids *[]uint) error {
	var children []platformdb.SysAuthority
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&children).Error; err != nil {
		return err
	}
	for _, child := range children {
		*ids = append(*ids, child.AuthorityId)
		if err := r.collectDescendants(ctx, child.AuthorityId, ids); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) CheckAuthorityAuth(ctx context.Context, adminID, targetID uint) error {
	if adminID == targetID {
		return nil
	}
	descendants, err := r.GetDescendantIDs(ctx, adminID)
	if err != nil {
		return err
	}
	for _, id := range descendants {
		if id == targetID {
			return nil
		}
	}
	return domain.ErrRoleNotFound
}

func isRoot(parentID *uint) bool {
	return parentID == nil || *parentID == 0
}
