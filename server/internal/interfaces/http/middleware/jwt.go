package middleware

import (
	"errors"
	"strconv"
	"time"

	platformauth "github.com/flipped-aurora/gin-vue-admin/server/internal/platform/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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
		token := utils.GetToken(c)
		if token == "" {
			response.Fail(c, 401, 7, "未登录或非法访问，请登录")
			c.Abort()
			return
		}
		if cfg.BlacklistCheck != nil && cfg.BlacklistCheck(token) {
			response.Fail(c, 401, 7, "您的账号已失效，请重新登录")
			utils.ClearToken(c)
			c.Abort()
			return
		}

		j := &utils.JWT{JWT: platformauth.NewJWT(platformauth.JWTConfig{SigningKey: cfg.SigningKey}), ConcurrencyControl: jwtConcurrency}
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.Fail(c, 401, 7, "登录已过期，请重新登录")
				utils.ClearToken(c)
				c.Abort()
				return
			}
			response.Fail(c, 401, 7, err.Error())
			utils.ClearToken(c)
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
			dr, _ := utils.ParseDuration(cfg.ExpiresTime)
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
			newToken, _ := j.CreateTokenByOldToken(token, *claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
			utils.SetToken(c, newToken, int(dr.Seconds()/60))
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
