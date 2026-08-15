package application

import (
	"context"
	"errors"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/user/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

type AuthorityChecker func(ctx context.Context, adminID, targetID uint) error

type Service struct {
	repo             domain.Repository
	hasher           platformauth.PasswordHasher
	authorityChecker AuthorityChecker
	strictAuth       bool
}

func NewService(repo domain.Repository, hasher platformauth.PasswordHasher, checker AuthorityChecker, strictAuth bool) *Service {
	if hasher == nil {
		hasher = platformauth.NewBcryptPasswordHasher()
	}
	return &Service{repo: repo, hasher: hasher, authorityChecker: checker, strictAuth: strictAuth}
}

// checkAuthorityScope 校验 targetID（角色 ID）是否在 adminID（调用者角色 ID）的层级作用域内。
// 与 role 模块约定一致：仅在 strictAuth 开启时强制执行层级校验；关闭时任意角色均可分配。
// checker 为 nil 时跳过（此时仓储通常也不可用，后续写入会失败）。
func (s *Service) checkAuthorityScope(ctx context.Context, adminID, targetID uint) error {
	if !s.strictAuth || s.authorityChecker == nil {
		return nil
	}
	if err := s.authorityChecker(ctx, adminID, targetID); err != nil {
		return apperrors.WithMessage(apperrors.Forbidden, "authority out of scope")
	}
	return nil
}

func (s *Service) Current(ctx context.Context) (CurrentUserResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return CurrentUserResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}

	if s.repo == nil {
		return CurrentUserResponse{User: userFromActor(actor)}, nil
	}

	user, err := s.repo.FindByID(ctx, actor.UserID)
	if err == nil {
		return CurrentUserResponse{User: userFromDomain(user)}, nil
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return CurrentUserResponse{User: userFromActor(actor)}, nil
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		return CurrentUserResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
	}
	return CurrentUserResponse{}, apperrors.New(apperrors.Internal, 0, "load current user failed", err)
}

func (s *Service) List(ctx context.Context, query ListUsersQuery) (ListUsersResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return ListUsersResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}

	page := pagination.Normalize(query.Page)
	if s.repo == nil {
		return emptyList(page), nil
	}

	result, err := s.repo.List(ctx, domain.ListQuery{
		Page:     page,
		Username: query.Username,
		NickName: query.NickName,
		Phone:    query.Phone,
		Email:    query.Email,
	})
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return emptyList(page), nil
	}
	if err != nil {
		return ListUsersResponse{}, apperrors.New(apperrors.Internal, 0, "list users failed", err)
	}

	users := make([]SourceUser, 0, len(result.List))
	for _, user := range result.List {
		users = append(users, userFromDomain(user))
	}
	return ListUsersResponse{
		List:     users,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

func (s *Service) ChangePassword(ctx context.Context, cmd ChangePasswordCommand) (ChangePasswordResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return ChangePasswordResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.OldPassword == "" || cmd.NewPassword == "" {
		return ChangePasswordResponse{}, apperrors.WithMessage(apperrors.Validation, "password is required")
	}
	if s.repo == nil {
		return ChangePasswordResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}

	oldHash, err := s.repo.FindPasswordHashByID(ctx, actor.UserID)
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return ChangePasswordResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		return ChangePasswordResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
	}
	if err != nil {
		return ChangePasswordResponse{}, apperrors.New(apperrors.Internal, 0, "load password failed", err)
	}
	if !s.hasher.Check(cmd.OldPassword, oldHash) {
		return ChangePasswordResponse{}, apperrors.WithMessage(apperrors.Validation, "old password is incorrect")
	}

	newHash, err := s.hasher.Hash(cmd.NewPassword)
	if err != nil {
		return ChangePasswordResponse{}, apperrors.New(apperrors.Internal, 0, "hash password failed", err)
	}
	if err := s.repo.UpdatePasswordHash(ctx, actor.UserID, newHash); err != nil {
		if errors.Is(err, domain.ErrRepositoryUnavailable) {
			return ChangePasswordResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
		}
		return ChangePasswordResponse{}, apperrors.New(apperrors.Internal, 0, "change password failed", err)
	}
	return ChangePasswordResponse{Changed: true}, nil
}

func (s *Service) UpdateProfile(ctx context.Context, cmd UpdateProfileCommand) (UpdateProfileResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if s.repo == nil {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	user, err := s.repo.UpdateProfile(ctx, actor.UserID, domain.ProfilePatch{
		NickName:  cmd.NickName,
		HeaderImg: cmd.HeaderImg,
		Phone:     cmd.Phone,
		Email:     cmd.Email,
	})
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
	}
	if err != nil {
		return UpdateProfileResponse{}, apperrors.New(apperrors.Internal, 0, "update profile failed", err)
	}
	return UpdateProfileResponse{User: userFromDomain(user)}, nil
}

// UpdateByAdmin 管理员更新目标用户的资料/状态/角色（与 profile 自改不同，作用于指定用户）。
func (s *Service) UpdateByAdmin(ctx context.Context, userID uint, cmd AdminUpdateUserCommand) (UpdateProfileResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if userID == 0 {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Validation, "user id is required")
	}
	if s.repo == nil {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}

	target, err := s.repo.FindByID(ctx, userID)
	if errors.Is(err, domain.ErrUserNotFound) {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	if err != nil {
		return UpdateProfileResponse{}, apperrors.New(apperrors.Internal, 0, "load target user failed", err)
	}
	if err := s.checkAuthorityScope(ctx, actor.AuthorityID, target.AuthorityID); err != nil {
		return UpdateProfileResponse{}, err
	}

	enable := cmd.Enable
	if enable == 0 {
		enable = target.Enable
	}
	if err := s.repo.UpdateByAdmin(ctx, userID, domain.AdminUpdateInput{
		NickName:  cmd.NickName,
		HeaderImg: cmd.HeaderImg,
		Phone:     cmd.Phone,
		Email:     cmd.Email,
		Enable:    enable,
	}); err != nil {
		if errors.Is(err, domain.ErrRepositoryUnavailable) {
			return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
		}
		if errors.Is(err, domain.ErrUserNotFound) {
			return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
		}
		return UpdateProfileResponse{}, apperrors.New(apperrors.Internal, 0, "update user failed", err)
	}

	if len(cmd.AuthorityIDs) > 0 {
		for _, id := range cmd.AuthorityIDs {
			if err := s.checkAuthorityScope(ctx, actor.AuthorityID, id); err != nil {
				return UpdateProfileResponse{}, err
			}
		}
		if err := s.repo.SetAuthorities(ctx, domain.SetAuthoritiesInput{UserID: userID, AuthorityIDs: cmd.AuthorityIDs}); err != nil {
			if errors.Is(err, domain.ErrRepositoryUnavailable) {
				return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
			}
			if errors.Is(err, domain.ErrUserNotFound) {
				return UpdateProfileResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
			}
			return UpdateProfileResponse{}, apperrors.New(apperrors.Internal, 0, "set user authorities failed", err)
		}
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return UpdateProfileResponse{}, apperrors.New(apperrors.Internal, 0, "refresh user after update failed", err)
	}
	return UpdateProfileResponse{User: userFromDomain(user)}, nil
}

func emptyList(page pagination.Page) ListUsersResponse {
	return ListUsersResponse{
		List:     []SourceUser{},
		Total:    0,
		Page:     page.Page,
		PageSize: page.PageSize,
	}
}

func (s *Service) Create(ctx context.Context, cmd CreateUserCommand) (CreateUserResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return CreateUserResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.Username == "" || cmd.Password == "" {
		return CreateUserResponse{}, apperrors.WithMessage(apperrors.Validation, "username and password are required")
	}
	if cmd.AuthorityID == 0 {
		return CreateUserResponse{}, apperrors.WithMessage(apperrors.Validation, "authority is required")
	}
	if s.repo == nil {
		return CreateUserResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	if err := s.checkAuthorityScope(ctx, actor.AuthorityID, cmd.AuthorityID); err != nil {
		return CreateUserResponse{}, err
	}
	for _, id := range cmd.AuthorityIDs {
		if err := s.checkAuthorityScope(ctx, actor.AuthorityID, id); err != nil {
			return CreateUserResponse{}, err
		}
	}

	passwordHash, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return CreateUserResponse{}, apperrors.New(apperrors.Internal, 0, "hash password failed", err)
	}

	user, err := s.repo.Create(ctx, domain.CreateUserInput{
		Username:     cmd.Username,
		PasswordHash: passwordHash,
		NickName:     cmd.NickName,
		HeaderImg:    cmd.HeaderImg,
		AuthorityID:  cmd.AuthorityID,
		AuthorityIDs: cmd.AuthorityIDs,
		Enable:       cmd.Enable,
		Phone:        cmd.Phone,
		Email:        cmd.Email,
	})
	if err == domain.ErrRepositoryUnavailable {
		return CreateUserResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	if err != nil {
		return CreateUserResponse{}, apperrors.New(apperrors.Internal, 0, "create user failed", err)
	}
	return CreateUserResponse{User: userFromDomain(user)}, nil
}

func (s *Service) Delete(ctx context.Context, userID uint) (DeleteUserResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return DeleteUserResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if actor.UserID == userID {
		return DeleteUserResponse{}, apperrors.WithMessage(apperrors.Validation, "cannot delete yourself")
	}
	if s.repo == nil {
		return DeleteUserResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}

	target, err := s.repo.FindByID(ctx, userID)
	if errors.Is(err, domain.ErrUserNotFound) {
		return DeleteUserResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return DeleteUserResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	if err != nil {
		return DeleteUserResponse{}, apperrors.New(apperrors.Internal, 0, "load target user failed", err)
	}
	if err := s.checkAuthorityScope(ctx, actor.AuthorityID, target.AuthorityID); err != nil {
		return DeleteUserResponse{}, err
	}

	if err := s.repo.Delete(ctx, userID); err == domain.ErrRepositoryUnavailable {
		return DeleteUserResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	} else if err != nil {
		return DeleteUserResponse{}, apperrors.New(apperrors.Internal, 0, "delete user failed", err)
	}
	return DeleteUserResponse{DeletedID: userID}, nil
}

func (s *Service) ResetPassword(ctx context.Context, cmd ResetPasswordCommand) (ResetPasswordResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.UserID == 0 || cmd.Password == "" {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.Validation, "user id and password are required")
	}
	if s.repo == nil {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}

	target, err := s.repo.FindByID(ctx, cmd.UserID)
	if errors.Is(err, domain.ErrUserNotFound) {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	if err != nil {
		return ResetPasswordResponse{}, apperrors.New(apperrors.Internal, 0, "load target user failed", err)
	}
	if err := s.checkAuthorityScope(ctx, actor.AuthorityID, target.AuthorityID); err != nil {
		return ResetPasswordResponse{}, err
	}

	passwordHash, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return ResetPasswordResponse{}, apperrors.New(apperrors.Internal, 0, "hash password failed", err)
	}
	if err := s.repo.UpdatePasswordHash(ctx, cmd.UserID, passwordHash); err == domain.ErrRepositoryUnavailable {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	} else if err != nil {
		return ResetPasswordResponse{}, apperrors.New(apperrors.Internal, 0, "reset password failed", err)
	}
	return ResetPasswordResponse{Changed: true}, nil
}

func (s *Service) SetAuthorities(ctx context.Context, cmd SetAuthoritiesCommand) (SetAuthoritiesResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if len(cmd.AuthorityIDs) == 0 {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.Validation, "at least one authority is required")
	}
	if s.repo == nil {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}

	target, err := s.repo.FindByID(ctx, cmd.UserID)
	if errors.Is(err, domain.ErrUserNotFound) {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
	}
	if errors.Is(err, domain.ErrRepositoryUnavailable) {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	if err != nil {
		return SetAuthoritiesResponse{}, apperrors.New(apperrors.Internal, 0, "load target user failed", err)
	}
	if err := s.checkAuthorityScope(ctx, actor.AuthorityID, target.AuthorityID); err != nil {
		return SetAuthoritiesResponse{}, err
	}
	for _, id := range cmd.AuthorityIDs {
		if err := s.checkAuthorityScope(ctx, actor.AuthorityID, id); err != nil {
			return SetAuthoritiesResponse{}, err
		}
	}

	if err := s.repo.SetAuthorities(ctx, domain.SetAuthoritiesInput{
		UserID:       cmd.UserID,
		AuthorityIDs: cmd.AuthorityIDs,
	}); err == domain.ErrRepositoryUnavailable {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	} else if err == domain.ErrUserNotFound {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.NotFound, "user not found")
	} else if err != nil {
		return SetAuthoritiesResponse{}, apperrors.New(apperrors.Internal, 0, "set user authorities failed", err)
	}

	user, err := s.repo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return SetAuthoritiesResponse{}, apperrors.New(apperrors.Internal, 0, "refresh user after set authorities failed", err)
	}
	return SetAuthoritiesResponse{User: userFromDomain(user)}, nil
}

func userFromActor(actor platformauth.Actor) SourceUser {
	return SourceUser{
		ID:           actor.UserID,
		Username:     actor.Username,
		NickName:     actor.NickName,
		AuthorityID:  actor.AuthorityID,
		Authorities:  []AuthorityRef{{AuthorityId: actor.AuthorityID}},
		AuthorityIds: []uint{actor.AuthorityID},
		Source:       "token",
	}
}

func userFromDomain(user domain.User) SourceUser {
	authorityIDs := user.AuthorityIDs
	if len(authorityIDs) == 0 && user.AuthorityID != 0 {
		authorityIDs = []uint{user.AuthorityID}
	}
	authorities := make([]AuthorityRef, 0, len(authorityIDs))
	for _, id := range authorityIDs {
		authorities = append(authorities, AuthorityRef{AuthorityId: id})
	}
	return SourceUser{
		ID:           user.ID,
		UUID:         user.UUID,
		Username:     user.Username,
		UserName:     user.Username,
		NickName:     user.NickName,
		HeaderImg:    user.HeaderImg,
		AuthorityID:  user.AuthorityID,
		Authorities:  authorities,
		AuthorityIds: authorityIDs,
		Phone:        user.Phone,
		Email:        user.Email,
		Enable:       user.Enable,
		Source:       "database",
	}
}
