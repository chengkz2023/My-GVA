package utils

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
)

type JWT struct {
	*platformauth.JWT
}

var (
	TokenValid            = platformauth.TokenValid
	TokenExpired          = platformauth.TokenExpired
	TokenNotValidYet      = platformauth.TokenNotValidYet
	TokenMalformed        = platformauth.TokenMalformed
	TokenSignatureInvalid = platformauth.TokenSignatureInvalid
	TokenInvalid          = platformauth.TokenInvalid
)

func NewJWT() *JWT {
	return &JWT{platformauth.NewJWT(platformauth.JWTConfig{
		SigningKey:  global.GVA_CONFIG.JWT.SigningKey,
		ExpiresTime: global.GVA_CONFIG.JWT.ExpiresTime,
		BufferTime:  global.GVA_CONFIG.JWT.BufferTime,
		Issuer:      global.GVA_CONFIG.JWT.Issuer,
	})}
}

func NewJWTWithKey(key []byte) *JWT {
	return &JWT{platformauth.NewJWT(platformauth.JWTConfig{
		SigningKey: string(key),
	})}
}

func (j *JWT) CreateClaims(baseClaims request.BaseClaims) request.CustomClaims {
	pc := j.JWT.CreateClaims(platformauth.BaseClaims{
		UUID:        baseClaims.UUID,
		ID:          baseClaims.ID,
		Username:    baseClaims.Username,
		NickName:    baseClaims.NickName,
		AuthorityId: baseClaims.AuthorityId,
	})
	return request.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: pc.BufferTime,
		RegisteredClaims: pc.RegisteredClaims,
	}
}

func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	return j.JWT.CreateToken(platformauth.CustomClaims{
		BaseClaims:       platformauth.BaseClaims(claims.BaseClaims),
		BufferTime:       claims.BufferTime,
		RegisteredClaims: claims.RegisteredClaims,
	})
}

func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	claims, err := j.JWT.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &request.CustomClaims{
		BaseClaims: request.BaseClaims{
			UUID:        claims.UUID,
			ID:          claims.BaseClaims.ID,
			Username:    claims.Username,
			NickName:    claims.NickName,
			AuthorityId: claims.AuthorityId,
		},
		BufferTime:       claims.BufferTime,
		RegisteredClaims: claims.RegisteredClaims,
	}, nil
}

func (j *JWT) CreateTokenByOldToken(oldToken string, claims request.CustomClaims) (string, error) {
	v, err, _ := global.GVA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.JWT.CreateToken(platformauth.CustomClaims{
			BaseClaims:       platformauth.BaseClaims(claims.BaseClaims),
			BufferTime:       claims.BufferTime,
			RegisteredClaims: claims.RegisteredClaims,
		})
	})
	return v.(string), err
}
