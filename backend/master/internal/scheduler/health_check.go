package scheduler

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// ScheduleHealthCheck 标记长时间未心跳的节点。
func ScheduleHealthCheck(db *gorm.DB, logger *logrus.Logger, timeout time.Duration) *Task {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return newTask(timeout, timeout, func() {
		threshold := time.Now().Add(-timeout)
		res := db.Model(&model.Node{}).
			Where("last_seen_at IS NULL OR last_seen_at < ?", threshold).
			Update("status", "offline")
		if res.Error != nil && logger != nil {
			logger.WithError(res.Error).Error("health check failed")
		}
	})
}
