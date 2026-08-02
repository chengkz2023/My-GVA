package utils

import (
	platformauth "github.com/chengkz2023/My-GVA/server/internal/platform/auth"
	"github.com/gin-gonic/gin"
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
