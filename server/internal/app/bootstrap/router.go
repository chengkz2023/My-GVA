package bootstrap

import (
	"net/http"
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func Routers(log *zap.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.GinRecovery(log, true))
	if gin.Mode() == gin.DebugMode {
		engine.Use(gin.Logger())
	}

	engine.StaticFS(global.GVA_CONFIG.Local.StorePath, justFilesFilesystem{http.Dir(global.GVA_CONFIG.Local.StorePath)})

	publicGroup := engine.Group(global.GVA_CONFIG.System.RouterPrefix)
	privateGroup := engine.Group(global.GVA_CONFIG.System.RouterPrefix)
	privateGroup.Use(middleware.JWTAuthWithConfig(middleware.JWTConfig{
		ExpiresTime: global.GVA_CONFIG.JWT.ExpiresTime,
		BufferTime:  global.GVA_CONFIG.JWT.BufferTime,
		SigningKey:  global.GVA_CONFIG.JWT.SigningKey,
		BlacklistCheck: func(token string) bool {
			_, ok := global.BlackCache.Get(token)
			return ok
		},
	})).Use(middleware.CasbinHandler())

	publicGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	initBizRouter(privateGroup, publicGroup)

	global.GVA_ROUTERS = engine.Routes()
	log.Info("router register success")
	return engine
}
