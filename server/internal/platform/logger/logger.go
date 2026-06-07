package logger

import (
	"fmt"
	"os"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/internal/platform/logger/internal"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() *zap.Logger {
	if ok, _ := utils.PathExists(global.GVA_CONFIG.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", global.GVA_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GVA_CONFIG.Zap.Director, os.ModePerm)
	}
	levels := global.GVA_CONFIG.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}
	logger := zap.New(zapcore.NewTee(cores...))
	opts := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel)}
	if global.GVA_CONFIG.Zap.ShowLine {
		opts = append(opts, zap.AddCaller())
	}
	return logger.WithOptions(opts...)
}

func Current() *zap.Logger {
	return global.GVA_LOG
}
