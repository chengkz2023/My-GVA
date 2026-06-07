package utils

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetToken(c *gin.Context) string {
	token := platformauth.GetToken(c)
	if token == "" {
		return token
	}
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		return token
	}
	platformauth.SetToken(c, token, int(claims.ExpiresAt.Unix()-time.Now().Unix()))
	return token
}

func SetToken(c *gin.Context, token string, maxAge int) {
	platformauth.SetToken(c, token, maxAge)
}

func ClearToken(c *gin.Context) {
	platformauth.ClearToken(c)
}

func GetClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	token := platformauth.GetToken(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

func GetUserID(c *gin.Context) uint {
	if claims, ok := platformauth.GetClaimsFromContext(c); ok {
		return claims.BaseClaims.ID
	}
	if cl, err := GetClaims(c); err == nil {
		return cl.BaseClaims.ID
	}
	return 0
}

func GetUserUuid(c *gin.Context) uuid.UUID {
	if claims, ok := platformauth.GetClaimsFromContext(c); ok {
		return claims.UUID
	}
	if cl, err := GetClaims(c); err == nil {
		return cl.UUID
	}
	return uuid.UUID{}
}

func GetUserAuthorityId(c *gin.Context) uint {
	if claims, ok := platformauth.GetClaimsFromContext(c); ok {
		return claims.AuthorityId
	}
	if cl, err := GetClaims(c); err == nil {
		return cl.AuthorityId
	}
	return 0
}

func GetUserInfo(c *gin.Context) *systemReq.CustomClaims {
	if claims, ok := platformauth.GetClaimsFromContext(c); ok {
		return &systemReq.CustomClaims{
			BaseClaims: systemReq.BaseClaims{
				UUID:        claims.UUID,
				ID:          claims.BaseClaims.ID,
				Username:    claims.Username,
				NickName:    claims.NickName,
				AuthorityId: claims.AuthorityId,
			},
			BufferTime:       claims.BufferTime,
			RegisteredClaims: claims.RegisteredClaims,
		}
	}
	if cl, err := GetClaims(c); err == nil {
		return cl
	}
	return nil
}

func GetUserName(c *gin.Context) string {
	if claims, ok := platformauth.GetClaimsFromContext(c); ok {
		return claims.Username
	}
	if cl, err := GetClaims(c); err == nil {
		return cl.Username
	}
	return ""
}

func LoginToken(user system.Login) (token string, claims systemReq.CustomClaims, err error) {
	j := NewJWT()
	claims = j.CreateClaims(systemReq.BaseClaims{
		UUID:        user.GetUUID(),
		ID:          user.GetUserId(),
		NickName:    user.GetNickname(),
		Username:    user.GetUsername(),
		AuthorityId: user.GetAuthorityId(),
	})
	token, err = j.CreateToken(claims)
	return
}
