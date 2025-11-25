package service

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
	"github.com/iwoov/snell-master/pkg/utils"
)

// NodeService 管理节点生命周期。
type NodeService struct {
	repo         repository.NodeRepository
	instanceRepo repository.InstanceRepository
	logger       *logrus.Logger
}

// NewNodeService 构造函数。
func NewNodeService(repo repository.NodeRepository, instanceRepo repository.InstanceRepository, logger *logrus.Logger) *NodeService {
	return &NodeService{repo: repo, instanceRepo: instanceRepo, logger: logger}
}

// RegisterNode 创建新节点并返回 API Token。
func (s *NodeService) RegisterNode(name, endpoint, location, countryCode string) (*model.Node, error) {
	token, err := utils.GenerateAPIToken()
	if err != nil {
		return nil, err
	}
	node := &model.Node{
		Name:        name,
		Endpoint:    endpoint,
		Location:    location,
		CountryCode: countryCode,
		APIToken:    token,
		Status:      "offline",
	}
	if err := s.repo.Create(node); err != nil {
		return nil, err
	}
	return node, nil
}

// GetNodeList 返回所有节点。
func (s *NodeService) GetNodeList() ([]model.Node, error) {
	return s.repo.List()
}

// GetNodeByID 查询节点详情。
func (s *NodeService) GetNodeByID(id uint) (*model.Node, error) {
	return s.repo.GetByID(id)
}

// UpdateNode 更新基础信息。
func (s *NodeService) UpdateNode(id uint, updates map[string]interface{}) (*model.Node, error) {
	node, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if name, ok := updates["name"].(string); ok && name != "" {
		node.Name = name
	}
	if endpoint, ok := updates["endpoint"].(string); ok && endpoint != "" {
		node.Endpoint = endpoint
	}
	if location, ok := updates["location"].(string); ok {
		node.Location = location
	}
	if code, ok := updates["country_code"].(string); ok {
		node.CountryCode = code
	}
	if status, ok := updates["status"].(string); ok && status != "" {
		node.Status = status
	}
	if err := s.repo.Update(node); err != nil {
		return nil, err
	}
	return node, nil
}

// DeleteNode 删除节点，若仍存在实例则报错。
func (s *NodeService) DeleteNode(id uint) error {
	instances, err := s.instanceRepo.GetByNode(id)
	if err != nil {
		return err
	}
	if len(instances) > 0 {
		return fmt.Errorf("node has active instances")
	}
	return s.repo.Delete(id)
}

// UpdateHeartbeat 更新节点心跳和统计。
func (s *NodeService) UpdateHeartbeat(apiToken string, cpu, mem float64, instanceCount int, status, version string) error {
	node, err := s.repo.GetByToken(apiToken)
	if err != nil {
		return err
	}
	if status == "" {
		status = "online"
	}
	if err := s.repo.UpdateHeartbeat(node.ID, cpu, mem, instanceCount, status); err != nil {
		return err
	}
	record := &model.NodeHeartbeat{
		NodeID:        node.ID,
		Status:        status,
		CPUUsage:      cpu,
		MemoryUsage:   mem,
		InstanceCount: instanceCount,
		Version:       version,
		CreatedAt:     time.Now(),
	}
	return s.repo.SaveHeartbeat(record)
}

// GetNodeByToken 根据 API Token 获取节点。
func (s *NodeService) GetNodeByToken(token string) (*model.Node, error) {
	return s.repo.GetByToken(token)
}

// RegenerateToken 重新生成 API Token。
func (s *NodeService) RegenerateToken(id uint) (string, error) {
	node, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateAPIToken()
	if err != nil {
		return "", err
	}
	node.APIToken = token
	if err := s.repo.Update(node); err != nil {
		return "", err
	}
	return token, nil
}
