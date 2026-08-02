package internal

import (
	"fmt"
	"github.com/chengkz2023/My-GVA/server/config"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type Writer struct {
	config config.GeneralDB
	writer logger.Writer
	log    *zap.Logger
}

func NewWriter(config config.GeneralDB, log *zap.Logger) *Writer {
	return &Writer{config: config, log: log}
}

// Printf 格式化打印日志
func (c *Writer) Printf(message string, data ...any) {

	// 当有日志时候均需要输出到控制台
	fmt.Printf(message, data...)

	// 当开启了zap的情况，会打印到日志记录
	if c.config.LogZap && c.log != nil {
		switch c.config.LogLevel() {
		case logger.Silent:
			c.log.Debug(fmt.Sprintf(message, data...))
		case logger.Error:
			c.log.Error(fmt.Sprintf(message, data...))
		case logger.Warn:
			c.log.Warn(fmt.Sprintf(message, data...))
		case logger.Info:
			c.log.Info(fmt.Sprintf(message, data...))
		default:
			c.log.Info(fmt.Sprintf(message, data...))
		}
		return
	}
}
