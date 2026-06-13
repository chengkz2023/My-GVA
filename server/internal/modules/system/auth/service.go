package auth

import (
	"context"

	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	apperrors "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/errors"
	legacymodel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Me(ctx context.Context) (MeResponse, error) {
	actor, ok := platformauth.ActorFromContext(ctx)
	if !ok {
		return MeResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "missing actor")
	}
	return MeResponse{
		Actor: ActorResponse{
			UserID:      actor.UserID,
			AuthorityID: actor.AuthorityID,
			Username:    actor.Username,
			NickName:    actor.NickName,
		},
	}, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (LoginResponse, error) {
	if s.db == nil {
		return LoginResponse{}, apperrors.WithMessage(apperrors.Internal, "database unavailable")
	}
	var user legacymodel.SysUser
	if err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return LoginResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "invalid username or password")
	}
	if !utils.BcryptCheck(password, user.Password) {
		return LoginResponse{}, apperrors.WithMessage(apperrors.Unauthorized, "invalid username or password")
	}
	if user.Enable != 1 {
		return LoginResponse{}, apperrors.WithMessage(apperrors.Forbidden, "account is disabled")
	}

	j := utils.NewJWT()
	claims := j.CreateClaims(request.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		return LoginResponse{}, apperrors.New(apperrors.Internal, 0, "create token failed", err)
	}
	return LoginResponse{
		User:  userToDTO(user),
		Token: token,
	}, nil
}

func userToDTO(user legacymodel.SysUser) UserDTO {
	return UserDTO{
		ID:          user.ID,
		UUID:        user.UUID.String(),
		Username:    user.Username,
		NickName:    user.NickName,
		AuthorityID: user.AuthorityId,
		HeaderImg:   user.HeaderImg,
	}
}
