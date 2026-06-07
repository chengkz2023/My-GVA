package auth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	TokenValid            = errors.New("unknown token error")
	TokenExpired          = errors.New("token expired")
	TokenNotValidYet      = errors.New("token not active yet")
	TokenMalformed        = errors.New("malformed token")
	TokenSignatureInvalid = errors.New("invalid token signature")
	TokenInvalid          = errors.New("invalid token")
)

type JWTConfig struct {
	SigningKey  string
	ExpiresTime string
	BufferTime  string
	Issuer      string
}

type JWT struct {
	SigningKey []byte
	Issuer     string
	Expires    time.Duration
	Buffer     time.Duration
}

func NewJWT(cfg JWTConfig) *JWT {
	expires, _ := time.ParseDuration(cfg.ExpiresTime)
	buffer, _ := time.ParseDuration(cfg.BufferTime)
	return &JWT{
		SigningKey: []byte(cfg.SigningKey),
		Issuer:     cfg.Issuer,
		Expires:    expires,
		Buffer:     buffer,
	}
}

func (j *JWT) CreateClaims(baseClaims BaseClaims) CustomClaims {
	return CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(j.Buffer / time.Second),
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"GVA"},
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.Expires)),
			Issuer:    j.Issuer,
		},
	}
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, TokenExpired
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, TokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, TokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, TokenNotValidYet
		default:
			return nil, TokenInvalid
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, TokenValid
}
