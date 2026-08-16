package audit

import (
	"context"

	platformdb "github.com/chengkz2023/My-GVA/server/internal/platform/database"
	"gorm.io/gorm"
)

// MySQLSink 审计默认实现：写入 sys_operation_records 表。
type MySQLSink struct {
	db *gorm.DB
}

func NewMySQLSink(db *gorm.DB) *MySQLSink {
	return &MySQLSink{db: db}
}

func (s *MySQLSink) Record(ctx context.Context, entry Entry) error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.WithContext(ctx).Create(&platformdb.SysOperationRecord{
		Ip:           entry.IP,
		Method:       entry.Method,
		Path:         entry.Path,
		Status:       entry.Status,
		Latency:      entry.Latency,
		Agent:        entry.Agent,
		ErrorMessage: entry.ErrorMessage,
		UserID:       entry.UserID,
	}).Error
}
