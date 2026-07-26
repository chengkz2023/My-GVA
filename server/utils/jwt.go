package utils

import (
	"golang.org/x/sync/singleflight"

	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
)

type JWT struct {
	*platformauth.JWT
	ConcurrencyControl *singleflight.Group
}

var (
	TokenValid            = platformauth.TokenValid
	TokenExpired          = platformauth.TokenExpired
	TokenNotValidYet      = platformauth.TokenNotValidYet
	TokenMalformed        = platformauth.TokenMalformed
	TokenSignatureInvalid = platformauth.TokenSignatureInvalid
	TokenInvalid          = platformauth.TokenInvalid
)

func NewJWTWithKey(key []byte) *JWT {
	return &JWT{JWT: platformauth.NewJWT(platformauth.JWTConfig{
		SigningKey: string(key),
	})}
}

func (j *JWT) CreateClaims(baseClaims platformauth.BaseClaims) platformauth.CustomClaims {
	pc := j.JWT.CreateClaims(platformauth.BaseClaims{
		UUID:        baseClaims.UUID,
		ID:          baseClaims.ID,
		Username:    baseClaims.Username,
		NickName:    baseClaims.NickName,
		AuthorityId: baseClaims.AuthorityId,
	})
	return platformauth.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: pc.BufferTime,
		RegisteredClaims: pc.RegisteredClaims,
	}
}

func (j *JWT) CreateToken(claims platformauth.CustomClaims) (string, error) {
	return j.JWT.CreateToken(platformauth.CustomClaims{
		BaseClaims:       platformauth.BaseClaims(claims.BaseClaims),
		BufferTime:       claims.BufferTime,
		RegisteredClaims: claims.RegisteredClaims,
	})
}

func (j *JWT) ParseToken(tokenString string) (*platformauth.CustomClaims, error) {
	claims, err := j.JWT.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &platformauth.CustomClaims{
		BaseClaims: platformauth.BaseClaims{
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

func (j *JWT) CreateTokenByOldToken(oldToken string, claims platformauth.CustomClaims) (string, error) {
	v, err, _ := j.ConcurrencyControl.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.JWT.CreateToken(platformauth.CustomClaims{
			BaseClaims:       platformauth.BaseClaims(claims.BaseClaims),
			BufferTime:       claims.BufferTime,
			RegisteredClaims: claims.RegisteredClaims,
		})
	})
	return v.(string), err
}
