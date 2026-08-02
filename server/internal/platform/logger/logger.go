package logger

import (
	"fmt"
	"os"

	"github.com/chengkz2023/My-GVA/server/config"
	"github.com/chengkz2023/My-GVA/server/internal/platform/logger/internal"
	"github.com/chengkz2023/My-GVA/server/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.Zap) *zap.Logger {
	if ok, _ := utils.PathExists(cfg.Director); !ok {
		fmt.Printf("create %v directory\n", cfg.Director)
		_ = os.Mkdir(cfg.Director, os.ModePerm)
	}
	levels := cfg.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i], cfg)
		cores = append(cores, core)
	}
	logger := zap.New(zapcore.NewTee(cores...))
	opts := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel)}
	if cfg.ShowLine {
		opts = append(opts, zap.AddCaller())
	}
	return logger.WithOptions(opts...)
}
