package internal

import (
	"os"
	"time"

	"github.com/chengkz2023/My-GVA/server/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapCore struct {
	level  zapcore.Level
	zapCfg config.Zap
	zapcore.Core
}

func NewZapCore(level zapcore.Level, zapCfg config.Zap) *ZapCore {
	entity := &ZapCore{level: level, zapCfg: zapCfg}
	syncer := entity.WriteSyncer()
	levelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == level
	})
	entity.Core = zapcore.NewCore(zapCfg.Encoder(), syncer, levelEnabler)
	return entity
}

func (z *ZapCore) WriteSyncer(formats ...string) zapcore.WriteSyncer {
	cutter := NewCutter(
		z.zapCfg.Director,
		z.level.String(),
		z.zapCfg.RetentionDay,
		CutterWithLayout(time.DateOnly),
		CutterWithFormats(formats...),
	)
	if z.zapCfg.LogInConsole {
		multiSyncer := zapcore.NewMultiWriteSyncer(os.Stdout, cutter)
		return zapcore.AddSync(multiSyncer)
	}
	return zapcore.AddSync(cutter)
}

func (z *ZapCore) Enabled(level zapcore.Level) bool {
	return z.level == level
}

func (z *ZapCore) With(fields []zapcore.Field) zapcore.Core {
	return z.Core.With(fields)
}

func (z *ZapCore) Check(entry zapcore.Entry, check *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(entry.Level) {
		return check.AddCore(entry, z)
	}
	return check
}

// Write 将日志写入目标。带 business/folder/directory 字段的条目会被路由到
// 对应子目录：这里只在本次写入范围内选择 core，绝不修改共享的 z.Core，
// 避免并发日志下的数据竞争（旧实现会替换 z.Core）。
func (z *ZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	core := z.Core
	for i := 0; i < len(fields); i++ {
		if fields[i].Key == "business" || fields[i].Key == "folder" || fields[i].Key == "directory" {
			syncer := z.WriteSyncer(fields[i].String)
			core = zapcore.NewCore(z.zapCfg.Encoder(), syncer, z.level)
		}
	}
	return core.Write(entry, fields)
}

func (z *ZapCore) Sync() error {
	return z.Core.Sync()
}
