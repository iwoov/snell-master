package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
	templatetool "github.com/iwoov/snell-master/backend/master/internal/template"
	"github.com/iwoov/snell-master/pkg/utils"
)

// SubscribeService 处理订阅令牌。
type SubscribeService struct {
	repo         repository.SubscribeRepository
	templateRepo repository.TemplateRepository
	userRepo     repository.UserRepository
	nodeRepo     repository.NodeRepository
	instanceRepo repository.InstanceRepository
	logger       *logrus.Logger
}

// NewSubscribeService 构造函数。
func NewSubscribeService(repo repository.SubscribeRepository, templateRepo repository.TemplateRepository, userRepo repository.UserRepository, nodeRepo repository.NodeRepository, instanceRepo repository.InstanceRepository, logger *logrus.Logger) *SubscribeService {
	return &SubscribeService{
		repo:         repo,
		templateRepo: templateRepo,
		userRepo:     userRepo,
		nodeRepo:     nodeRepo,
		instanceRepo: instanceRepo,
		logger:       logger,
	}
}

// CreateToken 创建订阅令牌，如果已存在则返回旧令牌。
func (s *SubscribeService) CreateToken(userID uint, templateID *uint, expiresAt *time.Time) (*model.SubscribeToken, error) {
	if existing, err := s.repo.GetByUser(userID); err == nil {
		return existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	token, err := utils.GenerateSubscribeToken()
	if err != nil {
		return nil, err
	}
	sub := &model.SubscribeToken{UserID: userID, Token: token, TemplateID: templateID, ExpiresAt: expiresAt}
	if err := s.repo.Create(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

// RegenerateToken 重新生成订阅令牌。
func (s *SubscribeService) RegenerateToken(userID uint) (string, error) {
	sub, err := s.repo.GetByUser(userID)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateSubscribeToken()
	if err != nil {
		return "", err
	}
	sub.Token = token
	if err := s.repo.Update(sub); err != nil {
		return "", err
	}
	return token, nil
}

// GetTokenByUser 返回用户订阅令牌。
func (s *SubscribeService) GetTokenByUser(userID uint) (*model.SubscribeToken, error) {
	return s.repo.GetByUser(userID)
}

// DeleteToken 删除记录。
func (s *SubscribeService) DeleteToken(id uint) error {
	return s.repo.Delete(id)
}

// IncrementAccess 记录访问。
func (s *SubscribeService) IncrementAccess(id uint) error {
	return s.repo.IncrementAccess(id)
}

// ListTokens 返回全部令牌。
func (s *SubscribeService) ListTokens() ([]model.SubscribeToken, error) {
	return s.repo.List()
}

// GetByToken 根据令牌查询。
func (s *SubscribeService) GetByToken(token string) (*model.SubscribeToken, error) {
	return s.repo.GetByToken(token)
}

// ValidateToken 校验令牌有效性。
func (s *SubscribeService) ValidateToken(token string) (*model.SubscribeToken, error) {
	sub, err := s.repo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	if sub.ExpiresAt != nil && time.Now().After(*sub.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}
	return sub, nil
}

// GenerateSurgeConfig 根据订阅令牌生成 Surge 配置。
func (s *SubscribeService) GenerateSurgeConfig(token string) (string, error) {
	sub, err := s.ValidateToken(token)
	if err != nil {
		return "", err
	}
	user, err := s.userRepo.GetByID(sub.UserID)
	if err != nil {
		return "", err
	}
	if user.Status == 0 {
		return "", fmt.Errorf("user disabled")
	}
	nodes, err := s.userRepo.GetUserNodes(sub.UserID)
	if err != nil {
		return "", err
	}
	instances, err := s.instanceRepo.GetByUser(sub.UserID)
	if err != nil {
		return "", err
	}

	var tpl *model.Template
	if sub.TemplateID != nil {
		tpl, err = s.templateRepo.GetByID(*sub.TemplateID)
		if err != nil {
			return "", err
		}
	} else {
		tpl, err = s.templateRepo.GetDefault()
		if err != nil {
			return "", err
		}
	}

	generator := templatetool.NewSurgeGenerator(tpl.Content)
	content, err := generator.Generate(user, nodes, instances)
	if err != nil {
		return "", err
	}
	if err := s.repo.IncrementAccess(sub.ID); err != nil {
		s.logger.WithError(err).Warn("increment subscribe access failed")
	}
	return content, nil
}
