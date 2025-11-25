package service

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
)

// TrafficService 负责流量记录与查询。
type TrafficService struct {
	repo     repository.TrafficRepository
	userRepo repository.UserRepository
	logger   *logrus.Logger
}

// NewTrafficService 构造函数。
func NewTrafficService(repo repository.TrafficRepository, userRepo repository.UserRepository, logger *logrus.Logger) *TrafficService {
	return &TrafficService{repo: repo, userRepo: userRepo, logger: logger}
}

// RecordTraffic 保存记录并更新用户统计。
func (s *TrafficService) RecordTraffic(userID, instanceID, nodeID uint, upload, download int64, recordDate time.Time) error {
	record := &model.TrafficRecord{
		UserID:        userID,
		InstanceID:    instanceID,
		NodeID:        nodeID,
		UploadBytes:   upload,
		DownloadBytes: download,
		RecordDate:    recordDate,
	}
	if err := s.repo.Create(record); err != nil {
		return err
	}
	total := upload + download
	if err := s.userRepo.UpdateTraffic(userID, total); err != nil {
		return err
	}
	return nil
}

// GetUserTraffic 返回时间范围内的数据。
func (s *TrafficService) GetUserTraffic(userID uint, start, end *time.Time) ([]model.TrafficRecord, error) {
	return s.repo.ListByUser(userID, start, end)
}

// GetTrafficSummary 返回整体统计。
func (s *TrafficService) GetTrafficSummary() (TrafficSummary, error) {
	summary, err := s.repo.GetSummary()
	if err != nil {
		return TrafficSummary{}, err
	}
	return TrafficSummary{
		TotalBytes:    summary.TotalBytes,
		TodayBytes:    summary.TodayBytes,
		RecordCount:   summary.RecordCount,
		UserCount:     summary.UserCount,
		InstanceCount: summary.InstanceCount,
	}, nil
}

// GetUserRanking 返回排行。
func (s *TrafficService) GetUserRanking(limit int) ([]UserTraffic, error) {
	stats, err := s.repo.GetUserRanking(limit)
	if err != nil {
		return nil, err
	}
	res := make([]UserTraffic, 0, len(stats))
	for _, stat := range stats {
		res = append(res, UserTraffic{UserID: stat.UserID, TotalBytes: stat.TotalBytes})
	}
	return res, nil
}

// GetNodeTraffic 返回某节点统计。
func (s *TrafficService) GetNodeTraffic(nodeID uint) (*NodeTraffic, error) {
	stat, err := s.repo.GetNodeTraffic(nodeID)
	if err != nil {
		return nil, err
	}
	return &NodeTraffic{NodeID: stat.NodeID, TotalBytes: stat.TotalBytes}, nil
}

// GetTrafficTrend 返回最近 N 天趋势。
func (s *TrafficService) GetTrafficTrend(days int) ([]TrafficTrendPoint, error) {
	if days <= 0 {
		days = 7
	}
	end := time.Now()
	start := end.AddDate(0, 0, -days)
	points, err := s.repo.GetTrend(start, end)
	if err != nil {
		return nil, err
	}
	res := make([]TrafficTrendPoint, 0, len(points))
	for _, p := range points {
		res = append(res, TrafficTrendPoint{Date: p.RecordDate, TotalBytes: p.TotalBytes})
	}
	return res, nil
}
