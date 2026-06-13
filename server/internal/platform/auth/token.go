package auth

import (
	"net"

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
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
	}
}

func ClearToken(c *gin.Context) {
	if c == nil {
		return
	}
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}
	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("x-token", "", -1, "/", host, false, false)
	}
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
