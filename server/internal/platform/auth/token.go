package auth

import (
	"github.com/gin-gonic/gin"
)

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get("x-token")
	if token == "" {
		token, _ = c.Cookie("x-token")
	}
	return token
}

func SetToken(c *gin.Context, token string, maxAge int) {
	if c == nil {
		return
	}
	// httpOnly=true：前端通过 x-token 请求头携带令牌，cookie 仅供同源自动携带，
	// 不让 JS 读取，降低 XSS 窃取面。domain 留空以保持 host-only，避免 Domain=localhost 被浏览器拒绝。
	c.SetCookie("x-token", token, maxAge, "/", "", false, true)
}

func ClearToken(c *gin.Context) {
	if c == nil {
		return
	}
	c.SetCookie("x-token", "", -1, "/", "", false, true)
}

func GetClaimsFromContext(c *gin.Context) (*CustomClaims, bool) {
	if c == nil {
		return nil, false
	}
	claims, ok := c.Get("claims")
	if !ok {
		return nil, false
	}
	customClaims, ok := claims.(*CustomClaims)
	return customClaims, ok
}
