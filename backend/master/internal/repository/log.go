package repository

import (
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// LogRepository 记录操作日志。
type LogRepository interface {
	Create(log *model.OperationLog) error
	List(offset, limit int, filters map[string]interface{}) ([]model.OperationLog, int64, error)
}

type logRepository struct {
	db *gorm.DB
}

// NewLogRepository 返回实现。
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

func (r *logRepository) Create(log *model.OperationLog) error {
	return r.db.Create(log).Error
}

func (r *logRepository) List(offset, limit int, filters map[string]interface{}) ([]model.OperationLog, int64, error) {
	query := r.db.Model(&model.OperationLog{})
	if v, ok := filters["admin_id"]; ok {
		query = query.Where("admin_id = ?", v)
	}
	if v, ok := filters["action"]; ok {
		query = query.Where("action = ?", v)
	}
	if v, ok := filters["target_type"]; ok {
		query = query.Where("target_type = ?", v)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []model.OperationLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
