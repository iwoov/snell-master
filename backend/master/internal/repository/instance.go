package repository

import (
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// InstanceFilter 允许查询过滤。
type InstanceFilter struct {
	UserID *uint
	NodeID *uint
	Status string
}

// InstanceRepository 定义实例访问接口。
type InstanceRepository interface {
	Create(instance *model.SnellInstance) error
	GetByID(id uint) (*model.SnellInstance, error)
	List(filter InstanceFilter) ([]model.SnellInstance, error)
	Update(instance *model.SnellInstance) error
	Delete(id uint) error
	UpdateStatus(id uint, status string) error
	GetByNode(nodeID uint) ([]model.SnellInstance, error)
	GetByUser(userID uint) ([]model.SnellInstance, error)
	CheckPortConflict(nodeID uint, port int) (bool, error)
}

type instanceRepository struct {
	db *gorm.DB
}

// NewInstanceRepository 构造实例。
func NewInstanceRepository(db *gorm.DB) InstanceRepository {
	return &instanceRepository{db: db}
}

func (r *instanceRepository) Create(instance *model.SnellInstance) error {
	return r.db.Create(instance).Error
}

func (r *instanceRepository) GetByID(id uint) (*model.SnellInstance, error) {
	var inst model.SnellInstance
	if err := r.db.Preload("User").Preload("Node").First(&inst, id).Error; err != nil {
		return nil, err
	}
	return &inst, nil
}

func (r *instanceRepository) List(filter InstanceFilter) ([]model.SnellInstance, error) {
	query := r.db.Preload("User").Preload("Node")
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.NodeID != nil {
		query = query.Where("node_id = ?", *filter.NodeID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var instances []model.SnellInstance
	if err := query.Order("id DESC").Find(&instances).Error; err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *instanceRepository) Update(instance *model.SnellInstance) error {
	return r.db.Save(instance).Error
}

func (r *instanceRepository) Delete(id uint) error {
	return r.db.Delete(&model.SnellInstance{}, id).Error
}

func (r *instanceRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.SnellInstance{}).Where("id = ?", id).Update("status", status).Error
}

func (r *instanceRepository) GetByNode(nodeID uint) ([]model.SnellInstance, error) {
	var instances []model.SnellInstance
	if err := r.db.Where("node_id = ?", nodeID).Find(&instances).Error; err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *instanceRepository) GetByUser(userID uint) ([]model.SnellInstance, error) {
	var instances []model.SnellInstance
	if err := r.db.Where("user_id = ?", userID).Find(&instances).Error; err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *instanceRepository) CheckPortConflict(nodeID uint, port int) (bool, error) {
	var count int64
	if err := r.db.Model(&model.SnellInstance{}).Where("node_id = ? AND port = ?", nodeID, port).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
