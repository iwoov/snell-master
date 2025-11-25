package service

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
	"github.com/iwoov/snell-master/pkg/utils"
)

// InstanceService 管理 Snell 实例。
type InstanceService struct {
	repo      repository.InstanceRepository
	userRepo  repository.UserRepository
	nodeRepo  repository.NodeRepository
	adminRepo repository.AdminRepository
	logger    *logrus.Logger
}

// NewInstanceService 构造函数。
func NewInstanceService(repo repository.InstanceRepository, userRepo repository.UserRepository, nodeRepo repository.NodeRepository, adminRepo repository.AdminRepository, logger *logrus.Logger) *InstanceService {
	return &InstanceService{repo: repo, userRepo: userRepo, nodeRepo: nodeRepo, adminRepo: adminRepo, logger: logger}
}

// CreateInstance 创建实例并分配端口。
func (s *InstanceService) CreateInstance(userID, nodeID uint, version int, obfs string) (*model.SnellInstance, error) {
	// 检查用户是否存在，如果不存在可能是管理员
	_, userErr := s.userRepo.GetByID(userID)
	if userErr != nil {
		// 尝试查找管理员
		admin, adminErr := s.adminRepo.GetByID(userID)
		if adminErr != nil {
			return nil, fmt.Errorf("user not found: %w", userErr)
		}

		// 找到管理员，为其创建对应的 user 记录
		user := &model.User{
			ID:           admin.ID,
			Username:     admin.Username,
			PasswordHash: admin.PasswordHash,
			Email:        "",
			TrafficLimit: 1099511627776, // 默认 1TB
			Status:       1,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, fmt.Errorf("create user from admin: %w", err)
		}
		s.logger.Infof("Auto-created user record for admin: %s (ID: %d)", admin.Username, admin.ID)
	}

	// 检查节点
	if _, err := s.nodeRepo.GetByID(nodeID); err != nil {
		return nil, err
	}

	port := utils.AllocatePort(userID)
	for {
		conflict, err := s.repo.CheckPortConflict(nodeID, port)
		if err != nil {
			return nil, err
		}
		if !conflict {
			break
		}
		port++
		if port > 60000 {
			return nil, fmt.Errorf("no available port")
		}
	}

	psk, err := utils.GeneratePSK()
	if err != nil {
		return nil, err
	}

	inst := &model.SnellInstance{
		UserID:  userID,
		NodeID:  nodeID,
		Port:    port,
		PSK:     psk,
		Version: version,
		Obfs:    obfs,
		Status:  "running",
	}
	if err := s.repo.Create(inst); err != nil {
		return nil, err
	}
	return inst, nil
}

// GetInstanceList 返回过滤后的实例。
func (s *InstanceService) GetInstanceList(filter repository.InstanceFilter) ([]model.SnellInstance, error) {
	return s.repo.List(filter)
}

// GetInstanceByID 查询详情。
func (s *InstanceService) GetInstanceByID(id uint) (*model.SnellInstance, error) {
	return s.repo.GetByID(id)
}

// DeleteInstance 删除实例。
func (s *InstanceService) DeleteInstance(id uint) error {
	return s.repo.Delete(id)
}

// UpdateInstanceStatus 更新状态。
func (s *InstanceService) UpdateInstanceStatus(id uint, status string) error {
	return s.repo.UpdateStatus(id, status)
}

// GetInstancesByNode 返回节点实例。
func (s *InstanceService) GetInstancesByNode(nodeID uint) ([]model.SnellInstance, error) {
	return s.repo.GetByNode(nodeID)
}

// GetInstancesByUser 返回用户实例。
func (s *InstanceService) GetInstancesByUser(userID uint) ([]model.SnellInstance, error) {
	return s.repo.GetByUser(userID)
}

// RestartInstance 当前仅记录日志，未来可调用 Agent。
func (s *InstanceService) RestartInstance(id uint) error {
	inst, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if s.logger != nil {
		s.logger.WithFields(logrus.Fields{
			"instance_id": id,
			"node_id":     inst.NodeID,
			"user_id":     inst.UserID,
		}).Info("instance restart requested")
	}
	return nil
}
