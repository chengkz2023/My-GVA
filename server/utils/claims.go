package utils

import (
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetToken(c *gin.Context) string {
	return platformauth.GetToken(c)
}

func SetToken(c *gin.Context, token string, maxAge int) {
	platformauth.SetToken(c, token, maxAge)
}

func ClearToken(c *gin.Context) {
	platformauth.ClearToken(c)
}

func GetClaims(c *gin.Context) (*systemReq.CustomClaims, error) {
	token := platformauth.GetToken(c)
	j := NewJWTWithKey(nil)
	claims, err := j.ParseToken(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
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

