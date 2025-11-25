package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// NodeRepository 操作节点数据。
type NodeRepository interface {
	Create(node *model.Node) error
	GetByID(id uint) (*model.Node, error)
	GetByToken(token string) (*model.Node, error)
	List() ([]model.Node, error)
	Update(node *model.Node) error
	Delete(id uint) error
	UpdateHeartbeat(nodeID uint, cpu, mem float64, instances int, status string) error
	SaveHeartbeat(record *model.NodeHeartbeat) error
	GetOnlineNodes(within time.Duration) ([]model.Node, error)
}

type nodeRepository struct {
	db *gorm.DB
}

// NewNodeRepository 构建实现。
func NewNodeRepository(db *gorm.DB) NodeRepository {
	return &nodeRepository{db: db}
}

func (r *nodeRepository) Create(node *model.Node) error {
	return r.db.Create(node).Error
}

func (r *nodeRepository) GetByID(id uint) (*model.Node, error) {
	var node model.Node
	if err := r.db.Preload("Instances").First(&node, id).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *nodeRepository) GetByToken(token string) (*model.Node, error) {
	var node model.Node
	if err := r.db.Where("api_token = ?", token).First(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *nodeRepository) List() ([]model.Node, error) {
	var nodes []model.Node
	if err := r.db.Order("id DESC").Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}

func (r *nodeRepository) Update(node *model.Node) error {
	return r.db.Save(node).Error
}

func (r *nodeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Node{}, id).Error
}

func (r *nodeRepository) UpdateHeartbeat(nodeID uint, cpu, mem float64, instances int, status string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"cpu_usage":      cpu,
		"memory_usage":   mem,
		"instance_count": instances,
		"status":         status,
		"last_seen_at":   &now,
		"updated_at":     now,
	}
	return r.db.Model(&model.Node{}).Where("id = ?", nodeID).Updates(updates).Error
}

func (r *nodeRepository) SaveHeartbeat(record *model.NodeHeartbeat) error {
	return r.db.Create(record).Error
}

func (r *nodeRepository) GetOnlineNodes(within time.Duration) ([]model.Node, error) {
	if within <= 0 {
		within = 5 * time.Minute
	}
	threshold := time.Now().Add(-within)
	var nodes []model.Node
	if err := r.db.Where("last_seen_at >= ?", threshold).Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}
