package mysql

import (
	"context"
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/menu/domain"
	legacysystem "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) TreeByAuthority(ctx context.Context, authorityID uint) ([]domain.Menu, error) {
	if r == nil || r.db == nil {
		return []domain.Menu{}, domain.ErrRepositoryUnavailable
	}

	var authorityMenus []legacysystem.SysAuthorityMenu
	if err := r.db.WithContext(ctx).
		Where("sys_authority_authority_id = ?", authorityID).
		Find(&authorityMenus).Error; err != nil {
		return nil, err
	}

	var menuIDs []uint
	for _, am := range authorityMenus {
		id, err := strconv.ParseUint(am.MenuId, 10, 64)
		if err != nil {
			continue
		}
		menuIDs = append(menuIDs, uint(id))
	}
	if len(menuIDs) == 0 {
		return []domain.Menu{}, nil
	}

	var baseMenus []legacysystem.SysBaseMenu
	if err := r.db.WithContext(ctx).
		Where("id IN ?", menuIDs).
		Order("sort").
		Find(&baseMenus).Error; err != nil {
		return nil, err
	}

	treeMap := make(map[uint][]legacysystem.SysBaseMenu)
	for _, bm := range baseMenus {
		treeMap[bm.ParentId] = append(treeMap[bm.ParentId], bm)
	}

	domainMenus := buildTree(treeMap, 0)
	return keepScaffoldMenus(domainMenus), nil
}

func buildTree(treeMap map[uint][]legacysystem.SysBaseMenu, parentID uint) []domain.Menu {
	baseMenus := treeMap[parentID]
	domainMenus := make([]domain.Menu, 0, len(baseMenus))
	for _, bm := range baseMenus {
		m := domain.Menu{
			ID:        bm.ID,
			ParentID:  bm.ParentId,
			Path:      bm.Path,
			Name:      bm.Name,
			Hidden:    bm.Hidden,
			Component: bm.Component,
			Sort:      bm.Sort,
			Title:     bm.Meta.Title,
			Icon:      bm.Meta.Icon,
		}
		m.Children = buildTree(treeMap, bm.ID)
		domainMenus = append(domainMenus, m)
	}
	return domainMenus
}

func (r *Repository) AssignMenus(ctx context.Context, authorityID uint, menuIDs []uint) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}

	menus := make([]legacysystem.SysBaseMenu, 0, len(menuIDs))
	for _, id := range menuIDs {
		menus = append(menus, legacysystem.SysBaseMenu{})
		menus[len(menus)-1].ID = id
	}

	var authority legacysystem.SysAuthority
	authority.AuthorityId = authorityID
	authority.SysBaseMenus = menus

	return r.db.WithContext(ctx).Model(&authority).Association("SysBaseMenus").Replace(menus)
}

func (r *Repository) Save(ctx context.Context, input domain.SaveMenuInput) (uint, error) {
	if r == nil || r.db == nil {
		return 0, domain.ErrRepositoryUnavailable
	}

	updating := input.ID != 0

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Name uniqueness check
		if updating {
			var old legacysystem.SysBaseMenu
			if err := tx.Where("id = ?", input.ID).First(&old).Error; err != nil {
				return domain.ErrMenuNotFound
			}
			if old.Name != input.Name {
				var count int64
				tx.Model(&legacysystem.SysBaseMenu{}).Where("id <> ? AND name = ?", input.ID, input.Name).Count(&count)
				if count > 0 {
					return domain.ErrMenuNameDuplicate
				}
			}
		} else {
			var count int64
			tx.Model(&legacysystem.SysBaseMenu{}).Where("name = ?", input.Name).Count(&count)
			if count > 0 {
				return domain.ErrMenuNameDuplicate
			}
		}

		// Parent check
		if input.ParentID != 0 {
			var parent legacysystem.SysBaseMenu
			if err := tx.Where("id = ?", input.ParentID).First(&parent).Error; err != nil {
				return domain.ErrParentNotFound
			}
			// Leaf-to-branch transition
			var childCount int64
			tx.Model(&legacysystem.SysBaseMenu{}).Where("parent_id = ?", input.ParentID).Count(&childCount)
			if childCount == 0 && !updating {
				var defaultCount int64
				tx.Model(&legacysystem.SysAuthority{}).Where("default_router = ?", parent.Name).Count(&defaultCount)
				if defaultCount > 0 {
					return domain.ErrParentIsDefaultRouter
				}
				tx.Where("sys_base_menu_id = ?", input.ParentID).Delete(&legacysystem.SysAuthorityMenu{})
			}
		}

		var menu legacysystem.SysBaseMenu
		if updating {
			menu.ID = input.ID
		}
		menu.ParentId = input.ParentID
		menu.Path = input.Path
		menu.Name = input.Name
		menu.Hidden = input.Hidden
		menu.Component = input.Component
		menu.Sort = input.Sort
		menu.Meta = legacysystem.Meta{
			Title:          input.Title,
			Icon:           input.Icon,
			ActiveName:     input.ActiveName,
			KeepAlive:      input.KeepAlive,
			DefaultMenu:    input.DefaultMenu,
			CloseTab:       input.CloseTab,
			TransitionType: input.TransitionType,
		}

		if updating {
			if err := tx.Save(&menu).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Create(&menu).Error; err != nil {
				return err
			}
		}

		// Rebuild parameters
		tx.Unscoped().Where("sys_base_menu_id = ?", menu.ID).Delete(&legacysystem.SysBaseMenuParameter{})
		for _, p := range input.Parameters {
			if err := tx.Create(&legacysystem.SysBaseMenuParameter{
				SysBaseMenuID: menu.ID,
				Type:          p.Type,
				Key:           p.Key,
				Value:         p.Value,
			}).Error; err != nil {
				return err
			}
		}

		// Rebuild buttons
		tx.Unscoped().Where("sys_base_menu_id = ?", menu.ID).Delete(&legacysystem.SysBaseMenuBtn{})
		for _, b := range input.Buttons {
			if err := tx.Create(&legacysystem.SysBaseMenuBtn{
				SysBaseMenuID: menu.ID,
				Name:          b.Name,
				Desc:          b.Desc,
			}).Error; err != nil {
				return err
			}
		}

		input.ID = menu.ID
		return nil
	})
	return input.ID, err
}

func (r *Repository) Delete(ctx context.Context, id uint) error {
	if r == nil || r.db == nil {
		return domain.ErrRepositoryUnavailable
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var menu legacysystem.SysBaseMenu
		if err := tx.Where("id = ?", id).First(&menu).Error; err != nil {
			return domain.ErrMenuNotFound
		}

		var childCount int64
		tx.Model(&legacysystem.SysBaseMenu{}).Where("parent_id = ?", id).Count(&childCount)
		if childCount > 0 {
			return domain.ErrMenuHasChildren
		}

		var defaultCount int64
		tx.Model(&legacysystem.SysAuthority{}).Where("default_router = ?", menu.Name).Count(&defaultCount)
		if defaultCount > 0 {
			return domain.ErrMenuIsDefaultRouter
		}

		tx.Delete(&legacysystem.SysBaseMenu{}, "id = ?", id)
		tx.Delete(&legacysystem.SysBaseMenuParameter{}, "sys_base_menu_id = ?", id)
		tx.Delete(&legacysystem.SysBaseMenuBtn{}, "sys_base_menu_id = ?", id)
		tx.Delete(&legacysystem.SysAuthorityBtn{}, "sys_menu_id = ?", id)
		tx.Delete(&legacysystem.SysAuthorityMenu{}, "sys_base_menu_id = ?", id)
		return nil
	})
}

func (r *Repository) FindByID(ctx context.Context, id uint) (domain.MenuDetail, error) {
	if r == nil || r.db == nil {
		return domain.MenuDetail{}, domain.ErrRepositoryUnavailable
	}

	var menu legacysystem.SysBaseMenu
	err := r.db.WithContext(ctx).
		Preload("Parameters").
		Preload("MenuBtn").
		Where("id = ?", id).
		First(&menu).Error
	if err != nil {
		return domain.MenuDetail{}, domain.ErrMenuNotFound
	}

	detail := domain.MenuDetail{
		Menu: domain.Menu{
			ID:        menu.ID,
			ParentID:  menu.ParentId,
			Path:      menu.Path,
			Name:      menu.Name,
			Hidden:    menu.Hidden,
			Component: menu.Component,
			Sort:      menu.Sort,
			Title:     menu.Meta.Title,
			Icon:      menu.Meta.Icon,
		},
		Meta: domain.MenuMeta{
			ActiveName:     menu.Meta.ActiveName,
			KeepAlive:      menu.Meta.KeepAlive,
			DefaultMenu:    menu.Meta.DefaultMenu,
			CloseTab:       menu.Meta.CloseTab,
			TransitionType: menu.Meta.TransitionType,
		},
	}
	for _, p := range menu.Parameters {
		detail.Parameters = append(detail.Parameters, domain.MenuParameter{
			Type: p.Type, Key: p.Key, Value: p.Value,
		})
	}
	for _, b := range menu.MenuBtn {
		detail.Buttons = append(detail.Buttons, domain.MenuButton{
			Name: b.Name, Desc: b.Desc,
		})
	}
	return detail, nil
}

// keepScaffoldMenus retains only the "superAdmin" root node and its descendants.
// This preserves compatibility with the old keepScaffoldSysMenus behavior.
// TODO: Revisit when the frontend no longer depends on the "superAdmin" root convention.
func keepScaffoldMenus(menus []domain.Menu) []domain.Menu {
	for _, menu := range menus {
		if menu.Name == "superAdmin" {
			return []domain.Menu{menu}
		}
	}
	return []domain.Menu{}
}
