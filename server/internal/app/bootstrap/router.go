package bootstrap

import (
	"net/http"

	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	"github.com/chengkz2023/My-GVA/server/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

func Routers(c *container.Container) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.GinRecovery(c.Logger, true))
	if gin.Mode() == gin.DebugMode {
		engine.Use(gin.Logger())
	}

	publicGroup := engine.Group(c.Config.System.RouterPrefix)
	publicGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	c.Logger.Info("router register success")
	return engine
}
