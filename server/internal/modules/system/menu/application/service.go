package application

import (
	"context"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/menu/domain"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
)

type AuthorityChecker = func(ctx context.Context, adminID, targetID uint) error

type Service struct {
	repo             domain.Repository
	authorityChecker AuthorityChecker
}

func NewService(repo domain.Repository, checker AuthorityChecker) *Service {
	return &Service{repo: repo, authorityChecker: checker}
}

func (s *Service) Tree(ctx context.Context) (TreeResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return TreeResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return TreeResponse{List: []MenuResponse{}}, nil
	}

	menus, err := s.repo.TreeByAuthority(ctx, actor.AuthorityID)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return TreeResponse{List: []MenuResponse{}}, nil
	}
	if err != nil {
		return TreeResponse{}, apperrors.New(apperrors.Internal, 0, "list menus failed", err)
	}
	return TreeResponse{List: mapMenus(menus)}, nil
}

func (s *Service) AssignAuthority(ctx context.Context, authorityID uint, menuIDs []uint) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "menu repository unavailable")
	}
	if s.authorityChecker != nil {
		actor, _ := platformauth.ActorFromContext(ctx)
		if err := s.authorityChecker(ctx, actor.AuthorityID, authorityID); err != nil {
			return apperrors.WithMessage(apperrors.Forbidden, "role out of scope")
		}
	}

	if err := s.repo.AssignMenus(ctx, authorityID, menuIDs); errors.Is(err, domain.ErrRepositoryUnavailable) {
		return apperrors.WithMessage(apperrors.Internal, "menu repository unavailable")
	} else if err != nil {
		return apperrors.New(apperrors.Internal, 0, "assign menu authority failed", err)
	}
	return nil
}

func (s *Service) Save(ctx context.Context, req SaveMenuRequest) (SaveMenuResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return SaveMenuResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return SaveMenuResponse{}, apperrors.WithMessage(apperrors.Internal, "menu repository unavailable")
	}

	input := domain.SaveMenuInput{
		ID:             req.ID,
		ParentID:       req.ParentID,
		Path:           req.Path,
		Name:           req.Name,
		Hidden:         req.Hidden,
		Component:      req.Component,
		Sort:           req.Sort,
		Title:          req.Title,
		Icon:           req.Icon,
		ActiveName:     req.ActiveName,
		KeepAlive:      req.KeepAlive,
		DefaultMenu:    req.DefaultMenu,
		CloseTab:       req.CloseTab,
		TransitionType: req.TransitionType,
	}
	for _, p := range req.Parameters {
		input.Parameters = append(input.Parameters, domain.MenuParameter{Type: p.Type, Key: p.Key, Value: p.Value})
	}
	for _, b := range req.Buttons {
		input.Buttons = append(input.Buttons, domain.MenuButton{Name: b.Name, Desc: b.Desc})
	}

	id, err := s.repo.Save(ctx, input)
	if errors.Is(err, domain.ErrMenuNameDuplicate) {
		return SaveMenuResponse{}, apperrors.WithMessage(apperrors.Conflict, "菜单name重复")
	}
	if errors.Is(err, domain.ErrParentNotFound) {
		return SaveMenuResponse{}, apperrors.WithMessage(apperrors.Validation, "父菜单不存在")
	}
	if errors.Is(err, domain.ErrParentIsDefaultRouter) {
		return SaveMenuResponse{}, apperrors.WithMessage(apperrors.Conflict, "父菜单已被其他角色的首页占用")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return SaveMenuResponse{}, apperrors.WithMessage(apperrors.Internal, "menu repository unavailable")
	}
	if err != nil {
		return SaveMenuResponse{}, apperrors.New(apperrors.Internal, 0, "save menu failed", err)
	}
	return SaveMenuResponse{ID: id}, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "menu repository unavailable")
	}

	err := s.repo.Delete(ctx, id)
	if errors.Is(err, domain.ErrMenuNotFound) {
		return apperrors.WithMessage(apperrors.NotFound, "菜单不存在")
	}
	if errors.Is(err, domain.ErrMenuHasChildren) {
		return apperrors.WithMessage(apperrors.Conflict, "此菜单存在子菜单不可删除")
	}
	if errors.Is(err, domain.ErrMenuIsDefaultRouter) {
		return apperrors.WithMessage(apperrors.Conflict, "此菜单有角色正在作为首页不可删除")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return apperrors.WithMessage(apperrors.Internal, "menu repository unavailable")
	}
	if err != nil {
		return apperrors.New(apperrors.Internal, 0, "delete menu failed", err)
	}
	return nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (MenuDetailResponse, error) {
	if s.repo == nil {
		return MenuDetailResponse{}, apperrors.WithMessage(apperrors.Internal, "menu repository unavailable")
	}
	detail, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, domain.ErrMenuNotFound) {
		return MenuDetailResponse{}, apperrors.WithMessage(apperrors.NotFound, "菜单不存在")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return MenuDetailResponse{}, apperrors.WithMessage(apperrors.Internal, "menu repository unavailable")
	}
	if err != nil {
		return MenuDetailResponse{}, apperrors.New(apperrors.Internal, 0, "get menu failed", err)
	}
	resp := MenuDetailResponse{
		Menu: MenuResponse{
			ID:        detail.ID,
			ParentID:  detail.ParentID,
			Path:      detail.Path,
			Name:      detail.Name,
			Hidden:    detail.Hidden,
			Component: detail.Component,
			Sort:      detail.Sort,
			Title:     detail.Title,
			Icon:      detail.Icon,
		},
		Meta: MenuMetaResponse{
			ActiveName:     detail.Meta.ActiveName,
			KeepAlive:      detail.Meta.KeepAlive,
			DefaultMenu:    detail.Meta.DefaultMenu,
			CloseTab:       detail.Meta.CloseTab,
			TransitionType: detail.Meta.TransitionType,
		},
	}
	for _, p := range detail.Parameters {
		resp.Parameters = append(resp.Parameters, MenuParameterResponse{Type: p.Type, Key: p.Key, Value: p.Value})
	}
	for _, b := range detail.Buttons {
		resp.Buttons = append(resp.Buttons, MenuButtonResponse{Name: b.Name, Desc: b.Desc})
	}
	return resp, nil
}

func mapMenus(menus []domain.Menu) []MenuResponse {
	items := make([]MenuResponse, 0, len(menus))
	for _, menu := range menus {
		items = append(items, MenuResponse{
			ID:        menu.ID,
			ParentID:  menu.ParentID,
			Path:      menu.Path,
			Name:      menu.Name,
			Hidden:    menu.Hidden,
			Component: menu.Component,
			Sort:      menu.Sort,
			Title:     menu.Title,
			Icon:      menu.Icon,
			Children:  mapMenus(menu.Children),
		})
	}
	return items
}
