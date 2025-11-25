package repository

import (
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// TemplateRepository 管理 Surge 模板。
type TemplateRepository interface {
	Create(tpl *model.Template) error
	Update(tpl *model.Template) error
	Delete(id uint) error
	GetByID(id uint) (*model.Template, error)
	GetDefault() (*model.Template, error)
	List() ([]model.Template, error)
	SetDefault(id uint) error
}

type templateRepository struct {
	db *gorm.DB
}

// NewTemplateRepository 返回实现。
func NewTemplateRepository(db *gorm.DB) TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) Create(tpl *model.Template) error {
	return r.db.Create(tpl).Error
}

func (r *templateRepository) Update(tpl *model.Template) error {
	return r.db.Save(tpl).Error
}

func (r *templateRepository) Delete(id uint) error {
	return r.db.Delete(&model.Template{}, id).Error
}

func (r *templateRepository) GetByID(id uint) (*model.Template, error) {
	var tpl model.Template
	if err := r.db.First(&tpl, id).Error; err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (r *templateRepository) GetDefault() (*model.Template, error) {
	var tpl model.Template
	if err := r.db.Where("is_default = ?", true).First(&tpl).Error; err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (r *templateRepository) List() ([]model.Template, error) {
	var tpl []model.Template
	if err := r.db.Order("id DESC").Find(&tpl).Error; err != nil {
		return nil, err
	}
	return tpl, nil
}

func (r *templateRepository) SetDefault(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Template{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&model.Template{}).Where("id = ?", id).Update("is_default", true).Error
	})
}
