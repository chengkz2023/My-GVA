package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/app/container"
	v2http "github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/interfaces/http/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/modules"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() {
	c := Initialize()
	engine := Router(c)
	address := fmt.Sprintf(":%d", c.Config.System.Addr)

	fmt.Printf(`
BoyKing Admin V2 started, version: %s
Visit: http://127.0.0.1%s
`, global.Version, address)

	Serve(address, engine, 10*time.Minute, 10*time.Minute)
}

func Router(c *container.Container) *gin.Engine {
	if global.GVA_CONFIG.Redis.Enable {
		Redis()
	}
	engine := Routers(c.Logger)
	v2http.RegisterV2(engine, v2Config(c), modules.HTTPModules(c)...)
	return engine
}

func v2Config(c *container.Container) v2http.Config {
	return v2http.Config{
		JWT: middleware.JWTConfig{
			ExpiresTime: c.Config.JWT.ExpiresTime,
			BufferTime:  c.Config.JWT.BufferTime,
			SigningKey:  c.Config.JWT.SigningKey,
			BlacklistCheck: func(token string) bool {
				_, ok := global.BlackCache.Get(token)
				return ok
			},
		},
	}
}

func Serve(address string, router *gin.Engine, readTimeout, writeTimeout time.Duration) {
	srv := &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
			zap.L().Error("server start failed", zap.Error(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("shutting down web server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("web server shutdown failed", zap.Error(err))
	}

	zap.L().Info("web server stopped")
}
