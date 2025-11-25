package repository

import (
	"time"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"gorm.io/gorm"
)

// SystemConfigRepository 系统配置仓库接口
type SystemConfigRepository interface {
	Get(key string) (*model.SystemConfig, error)
	Set(key, value string) error
	GetAll() ([]model.SystemConfig, error)
	GetByKeys(keys []string) (map[string]string, error)
	BatchSet(configs map[string]string) error
}

type systemConfigRepository struct {
	db *gorm.DB
}

// NewSystemConfigRepository 创建系统配置仓库实例
func NewSystemConfigRepository(db *gorm.DB) SystemConfigRepository {
	return &systemConfigRepository{db: db}
}

// Get 根据 key 查询配置
func (r *systemConfigRepository) Get(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := r.db.Where("key = ?", key).First(&config).Error
	return &config, err
}

// Set 更新配置值
func (r *systemConfigRepository) Set(key, value string) error {
	return r.db.Model(&model.SystemConfig{}).
		Where("key = ?", key).
		Updates(map[string]interface{}{
			"value":      value,
			"updated_at": time.Now(),
		}).Error
}

// GetAll 查询所有配置
func (r *systemConfigRepository) GetAll() ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := r.db.Order("key ASC").Find(&configs).Error
	return configs, err
}

// GetByKeys 批量查询配置（返回 map）
func (r *systemConfigRepository) GetByKeys(keys []string) (map[string]string, error) {
	var configs []model.SystemConfig
	err := r.db.Where("key IN ?", keys).Find(&configs).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	return result, nil
}

// BatchSet 批量更新配置
func (r *systemConfigRepository) BatchSet(configs map[string]string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		for key, value := range configs {
			if err := tx.Model(&model.SystemConfig{}).
				Where("key = ?", key).
				Updates(map[string]interface{}{
					"value":      value,
					"updated_at": now,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
