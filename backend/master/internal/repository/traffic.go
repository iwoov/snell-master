package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// TrafficSummary 聚合系统整体流量。
type TrafficSummary struct {
	TotalBytes    int64
	TodayBytes    int64
	RecordCount   int64
	UserCount     int64
	InstanceCount int64
}

// UserTrafficStat 用户流量排行。
type UserTrafficStat struct {
	UserID     uint
	TotalBytes int64
}

// NodeTrafficStat 节点流量统计。
type NodeTrafficStat struct {
	NodeID     uint
	TotalBytes int64
}

// TrafficTrendPoint 用于折线图。
type TrafficTrendPoint struct {
	RecordDate time.Time
	TotalBytes int64
}

// TrafficRepository 处理流量记录。
type TrafficRepository interface {
	Create(record *model.TrafficRecord) error
	ListByUser(userID uint, start, end *time.Time) ([]model.TrafficRecord, error)
	GetSummary() (TrafficSummary, error)
	GetUserRanking(limit int) ([]UserTrafficStat, error)
	GetNodeTraffic(nodeID uint) (NodeTrafficStat, error)
	GetTrend(start, end time.Time) ([]TrafficTrendPoint, error)
}

type trafficRepository struct {
	db *gorm.DB
}

// NewTrafficRepository 构造实例。
func NewTrafficRepository(db *gorm.DB) TrafficRepository {
	return &trafficRepository{db: db}
}

func (r *trafficRepository) Create(record *model.TrafficRecord) error {
	record.TotalBytes = record.UploadBytes + record.DownloadBytes
	return r.db.Create(record).Error
}

func (r *trafficRepository) ListByUser(userID uint, start, end *time.Time) ([]model.TrafficRecord, error) {
	query := r.db.Where("user_id = ?", userID)
	if start != nil {
		query = query.Where("record_date >= ?", *start)
	}
	if end != nil {
		query = query.Where("record_date <= ?", *end)
	}
	var records []model.TrafficRecord
	if err := query.Order("record_date DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *trafficRepository) GetSummary() (TrafficSummary, error) {
	var summary TrafficSummary
	if err := r.db.Model(&model.TrafficRecord{}).Select("COALESCE(SUM(total_bytes),0)").Scan(&summary.TotalBytes).Error; err != nil {
		return summary, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	if err := r.db.Model(&model.TrafficRecord{}).Where("record_date >= ?", today).Select("COALESCE(SUM(total_bytes),0)").Scan(&summary.TodayBytes).Error; err != nil {
		return summary, err
	}

	if err := r.db.Model(&model.TrafficRecord{}).Count(&summary.RecordCount).Error; err != nil {
		return summary, err
	}
	if err := r.db.Model(&model.User{}).Count(&summary.UserCount).Error; err != nil {
		return summary, err
	}
	if err := r.db.Model(&model.SnellInstance{}).Count(&summary.InstanceCount).Error; err != nil {
		return summary, err
	}
	return summary, nil
}

func (r *trafficRepository) GetUserRanking(limit int) ([]UserTrafficStat, error) {
	if limit <= 0 {
		limit = 10
	}
	var stats []UserTrafficStat
	if err := r.db.Model(&model.TrafficRecord{}).
		Select("user_id, COALESCE(SUM(total_bytes),0) AS total_bytes").
		Group("user_id").Order("total_bytes DESC").Limit(limit).
		Scan(&stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *trafficRepository) GetNodeTraffic(nodeID uint) (NodeTrafficStat, error) {
	var stat NodeTrafficStat
	stat.NodeID = nodeID
	if err := r.db.Model(&model.TrafficRecord{}).
		Select("COALESCE(SUM(total_bytes),0)").Where("node_id = ?", nodeID).
		Scan(&stat.TotalBytes).Error; err != nil {
		return stat, err
	}
	return stat, nil
}

func (r *trafficRepository) GetTrend(start, end time.Time) ([]TrafficTrendPoint, error) {
	var points []TrafficTrendPoint
	if err := r.db.Model(&model.TrafficRecord{}).
		Select("DATE(record_date) AS record_date, COALESCE(SUM(total_bytes),0) AS total_bytes").
		Where("record_date BETWEEN ? AND ?", start, end).
		Group("DATE(record_date)").Order("DATE(record_date)").
		Scan(&points).Error; err != nil {
		return nil, err
	}
	return points, nil
}
