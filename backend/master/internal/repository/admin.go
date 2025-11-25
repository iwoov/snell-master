package repository

import (
	"github.com/iwoov/snell-master/backend/master/internal/model"
	"gorm.io/gorm"
)

// AdminRepository 定义管理员持久化操作。
type AdminRepository interface {
	Create(admin *model.Admin) error
	GetByID(id uint) (*model.Admin, error)
	GetByUsername(username string) (*model.Admin, error)
	Update(admin *model.Admin) error
	Delete(id uint) error
	List(offset, limit int) ([]model.Admin, int64, error)
	CountByRole(role int) (int64, error)
	GetAll() ([]model.Admin, error)
}

type adminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 返回接口实现。
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

func (r *adminRepository) GetByID(id uint) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) GetByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	if err := r.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *adminRepository) Update(admin *model.Admin) error {
	return r.db.Save(admin).Error
}

func (r *adminRepository) Delete(id uint) error {
	return r.db.Delete(&model.Admin{}, id).Error
}

func (r *adminRepository) List(offset, limit int) ([]model.Admin, int64, error) {
	query := r.db.Model(&model.Admin{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var admins []model.Admin
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&admins).Error; err != nil {
		return nil, 0, err
	}
	return admins, total, nil
}

func (r *adminRepository) CountByRole(role int) (int64, error) {
	var total int64
	if err := r.db.Model(&model.Admin{}).Where("role = ?", role).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *adminRepository) GetAll() ([]model.Admin, error) {
	var admins []model.Admin
	if err := r.db.Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}
