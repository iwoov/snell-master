package service

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
)

// TemplateService 管理 Surge 模板。
type TemplateService struct {
	repo   repository.TemplateRepository
	logger *logrus.Logger
}

// NewTemplateService 构造函数。
func NewTemplateService(repo repository.TemplateRepository, logger *logrus.Logger) *TemplateService {
	return &TemplateService{repo: repo, logger: logger}
}

// CreateTemplate 创建模板。
func (s *TemplateService) CreateTemplate(name, content, description string, isDefault bool) (*model.Template, error) {
	tpl := &model.Template{
		Name:        name,
		Content:     content,
		Description: description,
		IsDefault:   isDefault,
	}
	if err := s.repo.Create(tpl); err != nil {
		return nil, err
	}
	if isDefault {
		if err := s.repo.SetDefault(tpl.ID); err != nil {
			return nil, err
		}
	}
	return tpl, nil
}

// UpdateTemplate 编辑模板。
func (s *TemplateService) UpdateTemplate(id uint, updates map[string]interface{}) (*model.Template, error) {
	tpl, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if name, ok := updates["name"].(string); ok && name != "" {
		tpl.Name = name
	}
	if content, ok := updates["content"].(string); ok && content != "" {
		tpl.Content = content
	}
	if desc, ok := updates["description"].(string); ok {
		tpl.Description = desc
	}
	if err := s.repo.Update(tpl); err != nil {
		return nil, err
	}
	if isDefault, ok := updates["is_default"].(bool); ok && isDefault {
		if err := s.repo.SetDefault(tpl.ID); err != nil {
			return nil, err
		}
	}
	return tpl, nil
}

// DeleteTemplate 删除模板。
func (s *TemplateService) DeleteTemplate(id uint) error {
	tpl, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if tpl.IsDefault {
		return fmt.Errorf("cannot delete default template")
	}
	return s.repo.Delete(id)
}

// ListTemplates 返回全部模板。
func (s *TemplateService) ListTemplates() ([]model.Template, error) {
	return s.repo.List()
}

// GetDefaultTemplate 返回默认模板。
func (s *TemplateService) GetDefaultTemplate() (*model.Template, error) {
	return s.repo.GetDefault()
}

// SetDefaultTemplate 设置默认模板。
func (s *TemplateService) SetDefaultTemplate(id uint) error {
	return s.repo.SetDefault(id)
}
