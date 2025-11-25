package scheduler

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/iwoov/snell-master/backend/master/internal/repository"
)

// ScheduleDailyReset 每日 00:00 重置当天流量。
func ScheduleDailyReset(userRepo repository.UserRepository, logger *logrus.Logger) *Task {
	next := nextAt(0, 0)
	return newTask(next.Sub(time.Now()), 24*time.Hour, func() {
		if err := userRepo.ResetDailyTraffic(); err != nil && logger != nil {
			logger.WithError(err).Error("daily traffic reset failed")
		} else if logger != nil {
			logger.Info("daily traffic reset completed")
		}
	})
}

// ScheduleMonthlyReset 每月 1 日 00:10 重置月流量。
func ScheduleMonthlyReset(userRepo repository.UserRepository, logger *logrus.Logger) *Task {
	next := nextMonthAt(1, 0, 10)
	return newTask(next.Sub(time.Now()), 24*time.Hour, func() {
		now := time.Now()
		if now.Day() != 1 {
			return
		}
		if err := userRepo.ResetMonthlyTraffic(); err != nil && logger != nil {
			logger.WithError(err).Error("monthly traffic reset failed")
		} else if logger != nil {
			logger.Info("monthly traffic reset completed")
		}
	})
}

func nextAt(hour, min int) time.Time {
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, min, 0, 0, now.Location())
	if !target.After(now) {
		target = target.Add(24 * time.Hour)
	}
	return target
}

func nextMonthAt(day, hour, min int) time.Time {
	now := time.Now()
	target := time.Date(now.Year(), now.Month(), day, hour, min, 0, 0, now.Location())
	if !target.After(now) {
		target = target.AddDate(0, 1, 0)
		target = time.Date(target.Year(), target.Month(), day, hour, min, 0, 0, now.Location())
	}
	return target
}
