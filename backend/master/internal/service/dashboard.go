package service

import (
	"time"

	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// DashboardService 统计仪表盘。
type DashboardService struct {
	db *gorm.DB
}

// NewDashboardService 构造函数。
func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{db: db}
}

// GetStats 返回汇总信息。
func (s *DashboardService) GetStats() (*DashboardStats, error) {
	stats := &DashboardStats{}
	if err := s.db.Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.User{}).Where("status = ?", 1).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Node{}).Count(&stats.TotalNodes).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.Node{}).Where("status = ?", "online").Count(&stats.OnlineNodes).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.SnellInstance{}).Count(&stats.TotalInstances).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&model.TrafficRecord{}).Select("COALESCE(SUM(bytes_total),0)").Scan(&stats.TotalTraffic).Error; err != nil {
		return nil, err
	}
	today := time.Now().Truncate(24 * time.Hour)
	if err := s.db.Model(&model.TrafficRecord{}).Where("record_date >= ?", today).Select("COALESCE(SUM(bytes_total),0)").Scan(&stats.TodayTraffic).Error; err != nil {
		return nil, err
	}
	return stats, nil
}
