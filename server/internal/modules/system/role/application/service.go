package application

import (
	"context"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/role/domain"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/authz"
	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
)

var defaultCasbinPolicies = []authz.Policy{
	{Path: "/menu/getMenu", Method: "POST"},
	{Path: "/jwt/jsonInBlacklist", Method: "POST"},
	{Path: "/base/login", Method: "POST"},
	{Path: "/user/changePassword", Method: "POST"},
	{Path: "/user/setUserAuthority", Method: "POST"},
	{Path: "/user/getUserInfo", Method: "GET"},
	{Path: "/user/setSelfInfo", Method: "PUT"},
	{Path: "/fileUploadAndDownload/upload", Method: "POST"},
	{Path: "/sysDictionary/findSysDictionary", Method: "GET"},
}

type Service struct {
	repo           domain.Repository
	policyProvider authz.PolicyProvider
	policySyncer   authz.PolicySyncer
	strictAuth     bool
}

func NewService(repo domain.Repository, policyProvider authz.PolicyProvider, policySyncer authz.PolicySyncer, strictAuth bool) *Service {
	return &Service{repo: repo, policyProvider: policyProvider, policySyncer: policySyncer, strictAuth: strictAuth}
}

func (s *Service) FindByID(ctx context.Context, authorityID uint) (RoleResponse, error) {
	if s.repo == nil {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Internal, "role repository unavailable")
	}
	role, err := s.repo.FindByID(ctx, authorityID)
	if err != nil {
		return RoleResponse{}, err
	}
	return RoleResponse{
		AuthorityID: role.AuthorityID, AuthorityName: role.AuthorityName,
		ParentID: role.ParentID, DefaultRouter: role.DefaultRouter,
	}, nil
}

func (s *Service) Tree(ctx context.Context) (TreeResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return TreeResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return TreeResponse{List: []RoleResponse{}}, nil
	}

	roles, err := s.repo.Tree(ctx, actor.AuthorityID, s.strictAuth)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return TreeResponse{List: []RoleResponse{}}, nil
	}
	if err != nil {
		return TreeResponse{}, apperrors.New(apperrors.Internal, 0, "list roles failed", err)
	}
	return TreeResponse{List: mapRoles(roles)}, nil
}

func (s *Service) Create(ctx context.Context, input domain.SaveRoleInput) (RoleResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Internal, "role repository unavailable")
	}

	if s.strictAuth && (input.ParentID == nil || *input.ParentID == 0) {
		actor, _ := platformauth.ActorFromContext(ctx)
		actorID := actor.AuthorityID
		input.ParentID = &actorID
	}

	if err := s.repo.Save(ctx, input); errors.Is(err, domain.ErrRoleIDExists) {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Conflict, "存在相同角色id")
	} else if err != nil {
		return RoleResponse{}, apperrors.New(apperrors.Internal, 0, "create role failed", err)
	}

	if s.policySyncer != nil {
		_ = s.policySyncer.SyncPolicies(input.AuthorityID, defaultCasbinPolicies)
	}
	return RoleResponse{AuthorityID: input.AuthorityID, AuthorityName: input.AuthorityName, ParentID: input.ParentID, DefaultRouter: input.DefaultRouter}, nil
}

func (s *Service) Update(ctx context.Context, input domain.SaveRoleInput) (RoleResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Internal, "role repository unavailable")
	}
	if err := s.repo.Save(ctx, input); errors.Is(err, domain.ErrRoleNotFound) {
		return RoleResponse{}, apperrors.WithMessage(apperrors.NotFound, "角色不存在")
	} else if err != nil {
		return RoleResponse{}, apperrors.New(apperrors.Internal, 0, "update role failed", err)
	}
	return RoleResponse{AuthorityID: input.AuthorityID, AuthorityName: input.AuthorityName, ParentID: input.ParentID, DefaultRouter: input.DefaultRouter}, nil
}

func (s *Service) Delete(ctx context.Context, authorityID uint) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "role repository unavailable")
	}
	_, err := s.repo.Delete(ctx, authorityID)
	if errors.Is(err, domain.ErrRoleNotFound) {
		return apperrors.WithMessage(apperrors.NotFound, "角色不存在")
	}
	if errors.Is(err, domain.ErrRoleHasUsers) {
		return apperrors.WithMessage(apperrors.Conflict, "此角色有用户正在使用禁止删除")
	}
	if errors.Is(err, domain.ErrRoleHasChildren) {
		return apperrors.WithMessage(apperrors.Conflict, "此角色存在子角色不允许删除")
	}
	if err != nil {
		return apperrors.New(apperrors.Internal, 0, "delete role failed", err)
	}
	if s.policySyncer != nil {
		_ = s.policySyncer.SyncPolicies(authorityID, nil)
	}
	return nil
}

func (s *Service) Copy(ctx context.Context, oldAuthorityID uint, newRole domain.SaveRoleInput) (RoleResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Internal, "role repository unavailable")
	}

	if err := s.repo.Save(ctx, newRole); errors.Is(err, domain.ErrRoleIDExists) {
		return RoleResponse{}, apperrors.WithMessage(apperrors.Conflict, "存在相同角色id")
	} else if err != nil {
		return RoleResponse{}, apperrors.New(apperrors.Internal, 0, "copy role failed", err)
	}

	if err := s.repo.CopyMenusAndButtons(ctx, oldAuthorityID, newRole.AuthorityID); err != nil {
		_, _ = s.repo.Delete(ctx, newRole.AuthorityID)
		return RoleResponse{}, apperrors.New(apperrors.Internal, 0, "copy menus and buttons failed", err)
	}

	if s.policyProvider != nil && s.policySyncer != nil {
		policies, _ := s.policyProvider.Policies(oldAuthorityID)
		if len(policies) > 0 {
			if err := s.policySyncer.SyncPolicies(newRole.AuthorityID, policies); err != nil {
				_, _ = s.repo.Delete(ctx, newRole.AuthorityID)
				return RoleResponse{}, apperrors.New(apperrors.Internal, 0, "copy casbin policies failed", err)
			}
		}
	}
	return RoleResponse{AuthorityID: newRole.AuthorityID, AuthorityName: newRole.AuthorityName, ParentID: newRole.ParentID, DefaultRouter: newRole.DefaultRouter}, nil
}

func (s *Service) SetDataAuthority(ctx context.Context, input domain.DataAuthorityInput) error {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return apperrors.WithMessage(apperrors.Internal, "role repository unavailable")
	}
	if s.strictAuth {
		actor, _ := platformauth.ActorFromContext(ctx)
		for _, id := range input.DataAuthorityIDs {
			if err := s.repo.CheckAuthorityAuth(ctx, actor.AuthorityID, id); err != nil {
				return apperrors.WithMessage(apperrors.Forbidden, "role out of scope")
			}
		}
	}
	if err := s.repo.SetDataAuthority(ctx, input); err != nil {
		return apperrors.New(apperrors.Internal, 0, "set data authority failed", err)
	}
	return nil
}

func (s *Service) GetDataAuthorities(ctx context.Context, authorityID uint) ([]uint, error) {
	if s.repo == nil {
		return nil, nil
	}
	return s.repo.GetDataAuthorities(ctx, authorityID)
}

func mapRoles(roles []domain.Role) []RoleResponse {
	items := make([]RoleResponse, 0, len(roles))
	for _, role := range roles {
		items = append(items, RoleResponse{
			AuthorityID:   role.AuthorityID,
			AuthorityName: role.AuthorityName,
			ParentID:      role.ParentID,
			DefaultRouter: role.DefaultRouter,
			Children:      mapRoles(role.Children),
		})
	}
	return items
}
