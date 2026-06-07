package middleware

import (
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func CorsByRules(corsCfg config.CORS) gin.HandlerFunc {
	if corsCfg.Mode == "allow-all" {
		return Cors()
	}
	whitelist := corsCfg.Whitelist
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		var matched *config.CORSWhitelist
		for _, w := range whitelist {
			if origin == w.AllowOrigin {
				matched = &w
				break
			}
		}
		if matched != nil {
			c.Header("Access-Control-Allow-Origin", matched.AllowOrigin)
			c.Header("Access-Control-Allow-Headers", matched.AllowHeaders)
			c.Header("Access-Control-Allow-Methods", matched.AllowMethods)
			c.Header("Access-Control-Expose-Headers", matched.ExposeHeaders)
			if matched.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}
		if matched == nil && corsCfg.Mode == "strict-whitelist" && !(c.Request.Method == "GET" && c.Request.URL.Path == "/health") {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}
		c.Next()
	}
}
