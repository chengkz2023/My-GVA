package domain

import (
	"context"
	"errors"
)

var (
	ErrRepositoryUnavailable = errors.New("menu repository unavailable")
	ErrMenuNotFound          = errors.New("menu not found")
	ErrMenuNameDuplicate     = errors.New("menu name already exists")
	ErrMenuHasChildren       = errors.New("menu has children")
	ErrMenuIsDefaultRouter   = errors.New("menu is used as default router")
	ErrParentNotFound        = errors.New("parent menu not found")
	ErrParentIsDefaultRouter = errors.New("parent menu is used as default router by other roles")
)

type Repository interface {
	TreeByAuthority(ctx context.Context, authorityID uint) ([]Menu, error)
	All(ctx context.Context) ([]Menu, error)
	AssignMenus(ctx context.Context, authorityID uint, menuIDs []uint) error
	Save(ctx context.Context, input SaveMenuInput) (uint, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (MenuDetail, error)
}

type SaveMenuInput struct {
	ID           uint
	ParentID     uint
	Path         string
	Name         string
	Hidden       bool
	Component    string
	Sort         int
	Title        string
	Icon         string
	ActiveName   string
	KeepAlive    bool
	DefaultMenu  bool
	CloseTab     bool
	TransitionType string
	Parameters   []MenuParameter
	Buttons      []MenuButton
}

type MenuParameter struct {
	Type  string
	Key   string
	Value string
}

type MenuButton struct {
	Name string
	Desc string
}

type MenuDetail struct {
	Menu
	Meta       MenuMeta
	Parameters []MenuParameter
	Buttons    []MenuButton
}

type MenuMeta struct {
	ActiveName     string
	KeepAlive      bool
	DefaultMenu    bool
	CloseTab       bool
	TransitionType string
}
