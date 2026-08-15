package middleware

import (
	"errors"
	"strconv"
	"time"

	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/chengkz2023/My-GVA/server/internal/platform/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/sync/singleflight"
)

var jwtConcurrency = &singleflight.Group{}

type JWTConfig struct {
	ExpiresTime string
	BufferTime  string
	SigningKey  string
	// BlacklistCheck returns true if the token is blacklisted.
	BlacklistCheck func(token string) bool
}

func JWTAuthWithConfig(cfg JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := platformauth.GetToken(c)
		if token == "" {
			response.Fail(c, 401, 7, "未登录或非法访问，请登录")
			c.Abort()
			return
		}
		if cfg.BlacklistCheck != nil && cfg.BlacklistCheck(token) {
			response.Fail(c, 401, 7, "您的账号已失效，请重新登录")
			platformauth.ClearToken(c)
			c.Abort()
			return
		}

		j := platformauth.NewJWT(platformauth.JWTConfig{SigningKey: cfg.SigningKey})
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, platformauth.TokenExpired) {
				response.Fail(c, 401, 7, "登录已过期，请重新登录")
				platformauth.ClearToken(c)
				c.Abort()
				return
			}
			response.Fail(c, 401, 7, err.Error())
			platformauth.ClearToken(c)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Request = c.Request.WithContext(platformauth.ContextWithActor(c.Request.Context(), platformauth.Actor{
			UserID:      claims.BaseClaims.ID,
			AuthorityID: claims.AuthorityId,
			Username:    claims.Username,
			NickName:    claims.NickName,
		}))
		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			dr, _ := platformauth.ParseDuration(cfg.ExpiresTime)
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
			v, _, _ := jwtConcurrency.Do("JWT:"+token, func() (interface{}, error) {
				return j.CreateToken(*claims)
			})
			newToken := v.(string)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			platformauth.SetToken(c, newToken, int(dr.Seconds()))
		}

		c.Next()

		if newToken, exists := c.Get("new-token"); exists {
			c.Header("new-token", newToken.(string))
		}
		if newExpiresAt, exists := c.Get("new-expires-at"); exists {
			c.Header("new-expires-at", newExpiresAt.(string))
		}
	}
}
