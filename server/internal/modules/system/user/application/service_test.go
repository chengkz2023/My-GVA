package application

import (
	"context"
	"errors"
	"testing"

	"github.com/chengkz2023/My-GVA/server/internal/modules/system/user/domain"
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"github.com/chengkz2023/My-GVA/server/internal/platform/pagination"
)

func TestCurrentFallsBackToActorWhenRepositoryUnavailable(t *testing.T) {
	service := NewService(&fakeRepository{err: domain.ErrRepositoryUnavailable}, fakeHasher{}, nil, false, nil)

	got, err := service.Current(actorContext())
	if err != nil {
		t.Fatalf("Current() error = %v", err)
	}
	if got.User.Source != "token" || got.User.Username != "admin" || got.User.AuthorityID != 888 {
		t.Fatalf("user = %+v, want token admin authority 888", got.User)
	}
}

func TestCurrentUsesRepositoryUser(t *testing.T) {
	service := NewService(&fakeRepository{user: domain.User{
		ID:          1,
		UUID:        "user-uuid",
		Username:    "admin",
		NickName:    "Administrator",
		HeaderImg:   "avatar.png",
		AuthorityID: 888,
		Phone:       "10086",
		Email:       "admin@example.com",
		Enable:      1,
	}}, fakeHasher{}, nil, false, nil)

	got, err := service.Current(actorContext())
	if err != nil {
		t.Fatalf("Current() error = %v", err)
	}
	if got.User.Source != "database" || got.User.UUID != "user-uuid" || got.User.Email != "admin@example.com" {
		t.Fatalf("user = %+v, want database user", got.User)
	}
}

func TestCurrentMissingActor(t *testing.T) {
	_, err := NewService(nil, fakeHasher{}, nil, false, nil).Current(context.Background())
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Unauthorized {
		t.Fatalf("error = %v, want unauthorized app error", err)
	}
}

func TestCurrentUserNotFound(t *testing.T) {
	_, err := NewService(&fakeRepository{err: domain.ErrUserNotFound}, fakeHasher{}, nil, false, nil).Current(actorContext())
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.NotFound {
		t.Fatalf("error = %v, want not found app error", err)
	}
}

func TestListUsers(t *testing.T) {
	service := NewService(&fakeRepository{list: pagination.Result[domain.User]{
		List: []domain.User{{
			ID:          1,
			Username:    "admin",
			NickName:    "Administrator",
			AuthorityID: 888,
			Email:       "admin@example.com",
		}},
		Total:    1,
		Page:     2,
		PageSize: 20,
	}}, fakeHasher{}, nil, false, nil)

	got, err := service.List(actorContext(), ListUsersQuery{
		Page:     pagination.Page{Page: 2, PageSize: 20},
		Username: "admin",
	})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got.Total != 1 || got.Page != 2 || got.PageSize != 20 || len(got.List) != 1 {
		t.Fatalf("result = %+v, want one user on page 2", got)
	}
	if got.List[0].Source != "database" || got.List[0].Email != "admin@example.com" {
		t.Fatalf("user = %+v, want database source with email", got.List[0])
	}
}

func TestListUsersRepositoryUnavailable(t *testing.T) {
	service := NewService(&fakeRepository{err: domain.ErrRepositoryUnavailable}, fakeHasher{}, nil, false, nil)

	got, err := service.List(actorContext(), ListUsersQuery{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if got.Total != 0 || got.Page != 1 || got.PageSize != 10 || len(got.List) != 0 {
		t.Fatalf("result = %+v, want empty normalized page", got)
	}
}

func TestChangePassword(t *testing.T) {
	repo := &fakeRepository{passwordHash: "hash:old-password"}
	service := NewService(repo, fakeHasher{}, nil, false, nil)

	got, err := service.ChangePassword(actorContext(), ChangePasswordCommand{
		OldPassword: "old-password",
		NewPassword: "new-password",
	})
	if err != nil {
		t.Fatalf("ChangePassword() error = %v", err)
	}
	if !got.Changed {
		t.Fatalf("changed = false, want true")
	}
	if repo.updatedHash != "hash:new-password" {
		t.Fatalf("updated hash = %q, want hash:new-password", repo.updatedHash)
	}
}

func TestChangePasswordWrongOldPassword(t *testing.T) {
	_, err := NewService(&fakeRepository{passwordHash: "hash:old-password"}, fakeHasher{}, nil, false, nil).ChangePassword(actorContext(), ChangePasswordCommand{
		OldPassword: "wrong-password",
		NewPassword: "new-password",
	})
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Validation {
		t.Fatalf("error = %v, want validation app error", err)
	}
}

func TestChangePasswordRepositoryUnavailable(t *testing.T) {
	_, err := NewService(nil, fakeHasher{}, nil, false, nil).ChangePassword(actorContext(), ChangePasswordCommand{
		OldPassword: "old-password",
		NewPassword: "new-password",
	})
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Internal {
		t.Fatalf("error = %v, want internal app error", err)
	}
}

func TestUpdateProfile(t *testing.T) {
	repo := &fakeRepository{}
	service := NewService(repo, fakeHasher{}, nil, false, nil)

	got, err := service.UpdateProfile(actorContext(), UpdateProfileCommand{
		NickName:  "New Admin",
		HeaderImg: "avatar.png",
		Phone:     "10086",
		Email:     "admin@example.com",
	})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	if repo.updatedProfile.NickName != "New Admin" || repo.updatedID != 1 {
		t.Fatalf("updated profile = %+v id=%d, want New Admin id 1", repo.updatedProfile, repo.updatedID)
	}
	if got.User.NickName != "New Admin" || got.User.Source != "database" {
		t.Fatalf("user = %+v, want updated database user", got.User)
	}
}

func TestUpdateProfileRepositoryUnavailable(t *testing.T) {
	_, err := NewService(nil, fakeHasher{}, nil, false, nil).UpdateProfile(actorContext(), UpdateProfileCommand{NickName: "New Admin"})
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Internal {
		t.Fatalf("error = %v, want internal app error", err)
	}
}

func TestCreateUser(t *testing.T) {
	service := NewService(&fakeRepository{user: domain.User{ID: 2, UUID: "new", Username: "newuser", Enable: 1}}, fakeHasher{}, nil, false, nil)

	got, err := service.Create(actorContext(), CreateUserCommand{
		Username:     "newuser",
		Password:     "secret",
		NickName:     "New User",
		AuthorityID:  888,
		AuthorityIDs: []uint{888},
		Enable:       1,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if got.User.Username != "newuser" || got.User.Source != "database" {
		t.Fatalf("user = %+v, want newuser database", got.User)
	}
}

func TestCreateUserMissingActor(t *testing.T) {
	_, err := NewService(&fakeRepository{}, fakeHasher{}, nil, false, nil).Create(context.Background(), CreateUserCommand{Username: "x", Password: "y"})
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Unauthorized {
		t.Fatalf("error = %v, want unauthorized", err)
	}
}

func TestDeleteUser(t *testing.T) {
	got, err := NewService(&fakeRepository{}, fakeHasher{}, nil, false, nil).Delete(actorContext(), 99)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if got.DeletedID != 99 {
		t.Fatalf("deleted id = %d, want 99", got.DeletedID)
	}
}

func TestDeleteSelf(t *testing.T) {
	_, err := NewService(&fakeRepository{}, fakeHasher{}, nil, false, nil).Delete(actorContext(), 1)
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Validation {
		t.Fatalf("error = %v, want validation error", err)
	}
}

func TestResetPassword(t *testing.T) {
	repo := &fakeRepository{passwordHash: "old"}
	got, err := NewService(repo, fakeHasher{}, nil, false, nil).ResetPassword(actorContext(), ResetPasswordCommand{
		UserID:   2,
		Password: "newpass",
	})
	if err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	if !got.Changed {
		t.Fatal("changed = false, want true")
	}
	if repo.updatedHash != "hash:newpass" {
		t.Fatalf("updated hash = %q", repo.updatedHash)
	}
}

func TestUpdateByAdminNonStrictSkipsScopeCheck(t *testing.T) {
	deny := AuthorityChecker(func(ctx context.Context, adminID, targetID uint) error {
		return errors.New("denied")
	})
	service := NewService(&fakeRepository{user: domain.User{ID: 2, Username: "target", AuthorityID: 888}}, fakeHasher{}, deny, false, nil)

	got, err := service.UpdateByAdmin(actorContext(), 2, AdminUpdateUserCommand{
		NickName:     "Updated",
		Enable:       1,
		AuthorityIDs: []uint{999, 777},
	})
	if err != nil {
		t.Fatalf("UpdateByAdmin() error = %v, want success in non-strict mode", err)
	}
	if got.User.Username != "target" {
		t.Fatalf("user = %+v, want target", got.User)
	}
}

func TestUpdateByAdminStrictBlocksOutOfScopeRole(t *testing.T) {
	checker := AuthorityChecker(func(ctx context.Context, adminID, targetID uint) error {
		if adminID == targetID {
			return nil
		}
		return errors.New("out of scope")
	})
	service := NewService(&fakeRepository{user: domain.User{ID: 2, Username: "target", AuthorityID: 888}}, fakeHasher{}, checker, true, nil)

	_, err := service.UpdateByAdmin(actorContext(), 2, AdminUpdateUserCommand{
		NickName:     "Updated",
		Enable:       1,
		AuthorityIDs: []uint{999},
	})
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Forbidden {
		t.Fatalf("error = %v, want forbidden app error", err)
	}
}

func TestPasswordPolicyRejectsWeakPassword(t *testing.T) {
	service := NewService(&fakeRepository{}, fakeHasher{}, nil, false, platformauth.DefaultPasswordPolicy{})

	_, err := service.Create(actorContext(), CreateUserCommand{
		Username:     "newuser",
		Password:     "short",
		AuthorityID:  888,
		AuthorityIDs: []uint{888},
		Enable:       1,
	})
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) || appErr.Kind != apperrors.Validation {
		t.Fatalf("error = %v, want validation app error", err)
	}
}

func TestPasswordPolicyAllowsStrongPassword(t *testing.T) {
	service := NewService(&fakeRepository{user: domain.User{ID: 2, Username: "newuser", Enable: 1}}, fakeHasher{}, nil, false, platformauth.DefaultPasswordPolicy{})

	_, err := service.Create(actorContext(), CreateUserCommand{
		Username:     "newuser",
		Password:     "strongPass123",
		AuthorityID:  888,
		AuthorityIDs: []uint{888},
		Enable:       1,
	})
	if err != nil {
		t.Fatalf("Create() error = %v, want success", err)
	}
}

func actorContext() context.Context {
	return platformauth.ContextWithActor(context.Background(), platformauth.Actor{
		UserID:      1,
		AuthorityID: 888,
		Username:    "admin",
		NickName:    "Admin",
	})
}

type fakeRepository struct {
	user           domain.User
	list           pagination.Result[domain.User]
	passwordHash   string
	updatedHash    string
	updatedID      uint
	updatedProfile domain.ProfilePatch
	err            error
}

func (r fakeRepository) FindByID(ctx context.Context, id uint) (domain.User, error) {
	if r.err != nil {
		return domain.User{}, r.err
	}
	return r.user, nil
}

func (r fakeRepository) List(ctx context.Context, query domain.ListQuery) (pagination.Result[domain.User], error) {
	if r.err != nil {
		return pagination.Result[domain.User]{}, r.err
	}
	return r.list, nil
}

func (r *fakeRepository) FindPasswordHashByID(ctx context.Context, id uint) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return r.passwordHash, nil
}

func (r *fakeRepository) UpdatePasswordHash(ctx context.Context, id uint, passwordHash string) error {
	if r.err != nil {
		return r.err
	}
	r.updatedHash = passwordHash
	return nil
}

func (r *fakeRepository) UpdateProfile(ctx context.Context, id uint, profile domain.ProfilePatch) (domain.User, error) {
	if r.err != nil {
		return domain.User{}, r.err
	}
	r.updatedID = id
	r.updatedProfile = profile
	return domain.User{
		ID:          id,
		Username:    "admin",
		NickName:    profile.NickName,
		HeaderImg:   profile.HeaderImg,
		AuthorityID: 888,
		Phone:       profile.Phone,
		Email:       profile.Email,
		Enable:      1,
	}, nil
}

func (r *fakeRepository) UpdateByAdmin(ctx context.Context, id uint, patch domain.AdminUpdateInput) error {
	return r.err
}

func (r *fakeRepository) Create(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	if r.err != nil {
		return domain.User{}, r.err
	}
	return domain.User{
		ID:          2,
		UUID:        "new-uuid",
		Username:    input.Username,
		NickName:    input.NickName,
		AuthorityID: input.AuthorityID,
		Enable:      1,
	}, nil
}

func (r *fakeRepository) Delete(ctx context.Context, id uint) error {
	return r.err
}

func (r *fakeRepository) SetAuthorities(ctx context.Context, input domain.SetAuthoritiesInput) error {
	return r.err
}

type fakeHasher struct{}

func (fakeHasher) Hash(password string) (string, error) {
	return "hash:" + password, nil
}

func (fakeHasher) Check(password string, hash string) bool {
	return hash == "hash:"+password
}
