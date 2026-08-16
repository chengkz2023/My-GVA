package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chengkz2023/My-GVA/server/internal/app/container"
	apphttp "github.com/chengkz2023/My-GVA/server/internal/interfaces/http"
	"github.com/chengkz2023/My-GVA/server/internal/interfaces/http/middleware"
	"github.com/chengkz2023/My-GVA/server/internal/modules"
	"github.com/chengkz2023/My-GVA/server/internal/platform/audit"
	"github.com/chengkz2023/My-GVA/server/internal/platform/buildinfo"
	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() {
	c := Initialize()
	engine := Router(c)
	address := fmt.Sprintf(":%d", c.Config.System.Addr)

	fmt.Printf(`
	BoyKing Admin started, version: %s
	Visit: http://127.0.0.1%s
	`, buildinfo.Version, address)

	Serve(address, engine, 10*time.Minute, 10*time.Minute)
}

func Router(c *container.Container) *gin.Engine {
	if c.Config.Redis.Enable {
		c.Redis = InitRedis(c.Config.Redis, c.Logger)
	}
	engine := Routers(c)
	apphttp.RegisterRoutes(engine, apiConfig(c), modules.HTTPModules(c)...)
	c.Routes = engine.Routes()
	return engine
}

func apiConfig(c *container.Container) apphttp.Config {
	return apphttp.Config{
		AuditSink: audit.NewMySQLSink(c.DB),
		Logger:    c.Logger,
		Enforcer:  c.Enforcer,
		JWT: middleware.JWTConfig{
			ExpiresTime: c.Config.JWT.ExpiresTime,
			BufferTime:  c.Config.JWT.BufferTime,
			SigningKey:  c.Config.JWT.SigningKey,
			BlacklistCheck: func(token string) bool {
				return platformdb.IsJwtBlacklistedCached(c.DB, c.BlackCache, token)
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
