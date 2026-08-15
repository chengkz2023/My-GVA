package bootstrap

import (
	"net/http"
	"os"

	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	"github.com/chengkz2023/My-GVA/server/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		return nil, os.ErrPermission
	}
	return f, nil
}

func Routers(c *container.Container) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.GinRecovery(c.Logger, true))
	if gin.Mode() == gin.DebugMode {
		engine.Use(gin.Logger())
	}

	engine.StaticFS(c.Config.Local.StorePath, justFilesFilesystem{http.Dir(c.Config.Local.StorePath)})

	publicGroup := engine.Group(c.Config.System.RouterPrefix)
	publicGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	c.Logger.Info("router register success")
	return engine
}
