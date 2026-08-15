package bootstrap

import (
	"fmt"

	"github.com/chengkz2023/My-GVA/server/internal/app/task"
	"github.com/chengkz2023/My-GVA/server/internal/platform/timer"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func InitTimer(db *gorm.DB, t timer.Timer) {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		_, err := t.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(db)
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}
	}()
}
