package service

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
)

// LogService 记录/查询操作日志。
type LogService struct {
	repo   repository.LogRepository
	logger *logrus.Logger
}

// NewLogService 构造函数。
func NewLogService(repo repository.LogRepository, logger *logrus.Logger) *LogService {
	return &LogService{repo: repo, logger: logger}
}

// LogOperation 记录操作。
func (s *LogService) LogOperation(adminID uint, action, targetType string, targetID *uint, details, ip string) error {
	log := &model.OperationLog{
		AdminID:    adminID,
		Action:     action,
		TargetType: targetType,
		Details:    details,
		IPAddress:  ip,
		CreatedAt:  time.Now(),
	}
	if targetID != nil {
		log.TargetID = targetID
	}
	return s.repo.Create(log)
}

// GetLogs 分页查询。
func (s *LogService) GetLogs(page, pageSize int, filters map[string]interface{}) ([]model.OperationLog, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize, filters)
}
