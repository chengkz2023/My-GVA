package application

import (
	"context"
	"errors"

	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules/system/user/domain"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/pagination"
)

type AuthorityChecker func(ctx context.Context, adminID, targetID uint) error

type Service struct {
	repo             domain.Repository
	hasher           platformauth.PasswordHasher
	authorityChecker AuthorityChecker
}

func NewService(repo domain.Repository, hasher platformauth.PasswordHasher, checker AuthorityChecker) *Service {
	if hasher == nil {
		hasher = platformauth.NewBcryptPasswordHasher()
	}
	return &Service{repo: repo, hasher: hasher, authorityChecker: checker}
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

func emptyList(page pagination.Page) ListUsersResponse {
	return ListUsersResponse{
		List:     []SourceUser{},
		Total:    0,
		Page:     page.Page,
		PageSize: page.PageSize,
	}
}

func (s *Service) Create(ctx context.Context, cmd CreateUserCommand) (CreateUserResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return CreateUserResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.Username == "" || cmd.Password == "" {
		return CreateUserResponse{}, apperrors.WithMessage(apperrors.Validation, "username and password are required")
	}
	if s.repo == nil {
		return CreateUserResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
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

	if err := s.repo.Delete(ctx, userID); err == domain.ErrRepositoryUnavailable {
		return DeleteUserResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	} else if err != nil {
		return DeleteUserResponse{}, apperrors.New(apperrors.Internal, 0, "delete user failed", err)
	}
	return DeleteUserResponse{DeletedID: userID}, nil
}

func (s *Service) ResetPassword(ctx context.Context, cmd ResetPasswordCommand) (ResetPasswordResponse, error) {
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if cmd.UserID == 0 || cmd.Password == "" {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.Validation, "user id and password are required")
	}
	if s.repo == nil {
		return ResetPasswordResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
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
	if _, ok := platformauth.ActorFromContext(ctx); !ok {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	if len(cmd.AuthorityIDs) == 0 {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.Validation, "at least one authority is required")
	}
	if s.repo == nil {
		return SetAuthoritiesResponse{}, apperrors.WithMessage(apperrors.Internal, "user repository unavailable")
	}
	// TODO: Add hierarchical auth check (CheckAuthorityIDAuth) when role write operations are migrated

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
		ID:          actor.UserID,
		Username:    actor.Username,
		NickName:    actor.NickName,
		AuthorityID: actor.AuthorityID,
		Source:      "token",
	}
}

func userFromDomain(user domain.User) SourceUser {
	return SourceUser{
		ID:          user.ID,
		UUID:        user.UUID,
		Username:    user.Username,
		NickName:    user.NickName,
		HeaderImg:   user.HeaderImg,
		AuthorityID: user.AuthorityID,
		Phone:       user.Phone,
		Email:       user.Email,
		Enable:      user.Enable,
		Source:      "database",
	}
}
