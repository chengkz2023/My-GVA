// Package audit 定义操作审计的输出接缝。
// 默认实现为 MySQLSink（写入 sys_operation_records）；等保/合规项目可替换为 syslog、ELK 等实现。
package audit

import (
	"context"
	"time"
)

// Entry 一次操作审计记录。
type Entry struct {
	UserID       int
	IP           string
	Method       string
	Path         string
	Status       int
	Latency      time.Duration
	Agent        string
	ErrorMessage string
}

// Sink 审计输出接缝。
type Sink interface {
	Record(ctx context.Context, entry Entry) error
}
