package auth

import (
	"context"
	"time"

	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
	apperrors "github.com/chengkz2023/My-GVA/server/internal/platform/errors"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	jwt    *platformauth.JWT
	hasher platformauth.PasswordHasher
}

func NewService(db *gorm.DB, jwt *platformauth.JWT, hasher platformauth.PasswordHasher) *Service {
	return &Service{db: db, jwt: jwt, hasher: hasher}
}

func (s *Service) Me(ctx context.Context) (MeResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return MeResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	ui := UserInfoResponse{
		ID:       actor.UserID,
		NickName: actor.NickName,
		Authority: AuthorityInfo{
			DefaultRouter: "authority",
		},
		Authorities: []any{},
	}
	if s.db != nil {
		var user platformdb.SysUser
		if err := s.db.WithContext(ctx).
			Preload("Authority").
			Where("id = ?", actor.UserID).
			First(&user).Error; err == nil {
			ui.ID = user.ID
			ui.UUID = user.UUID.String()
			ui.NickName = user.NickName
			ui.HeaderImg = user.HeaderImg
			ui.Authority.DefaultRouter = user.Authority.DefaultRouter
			if ui.Authority.DefaultRouter == "" {
				ui.Authority.DefaultRouter = "authority"
			}
			if user.OriginSetting != nil {
				originSetting := map[string]any{}
				for k, v := range user.OriginSetting {
					originSetting[k] = v
				}
				ui.OriginSetting = originSetting
			}
		}
	}
	return MeResponse{UserInfo: ui}, nil
}

func (s *Service) Logout(ctx context.Context, token string, expiresAt time.Time) error {
	return platformdb.AddToJwtBlacklist(s.db, token, expiresAt)
}

func (s *Service) Login(ctx context.Context, username, password string) (LoginResponse, error) {
	if s.db == nil {
		return LoginResponse{}, apperrors.WithMessage(apperrors.Internal, "database unavailable")
	}
	var user platformdb.SysUser
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return LoginResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "invalid username or password")
	}
	if !s.hasher.Check(password, user.Password) {
		return LoginResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "invalid username or password")
	}
	if user.Enable != 1 {
		return LoginResponse{}, apperrors.WithMessage(apperrors.Forbidden, "account is disabled")
	}

	claims := s.jwt.CreateClaims(platformauth.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		Username:    user.Username,
		NickName:    user.NickName,
		AuthorityId: user.AuthorityId,
	})
	token, err := s.jwt.CreateToken(claims)
	if err != nil {
		return LoginResponse{}, apperrors.New(apperrors.Internal, 0, "create token failed", err)
	}
	return LoginResponse{
		User:  userToDTO(user),
		Token: token,
	}, nil
}

func userToDTO(user platformdb.SysUser) UserDTO {
	return UserDTO{
		ID:          user.ID,
		UUID:        user.UUID.String(),
		Username:    user.Username,
		NickName:    user.NickName,
		AuthorityID: user.AuthorityId,
		HeaderImg:   user.HeaderImg,
	}
}
