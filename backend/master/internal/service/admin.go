package service

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
	jwtutil "github.com/iwoov/snell-master/pkg/jwt"
	"github.com/iwoov/snell-master/pkg/utils"
)

// AdminService 封装管理员业务逻辑。
type AdminService struct {
	repo          repository.AdminRepository
	logger        *logrus.Logger
	jwtSecret     string
	jwtExpireHour int
}

// NewAdminService 创建实例。
func NewAdminService(repo repository.AdminRepository, logger *logrus.Logger, jwtSecret string, jwtExpireHour int) *AdminService {
	return &AdminService{repo: repo, logger: logger, jwtSecret: jwtSecret, jwtExpireHour: jwtExpireHour}
}

// Login 校验管理员凭据并返回 JWT。
func (s *AdminService) Login(username, password string) (string, *model.Admin, error) {
	admin, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, fmt.Errorf("admin not found")
		}
		return "", nil, err
	}

	if err := utils.CheckPassword(admin.PasswordHash, password); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	role := mapAdminRole(admin.Role)
	token, err := jwtutil.GenerateToken(admin.ID, admin.Username, role, s.jwtSecret, s.jwtExpireHour)
	if err != nil {
		return "", nil, err
	}
	return token, admin, nil
}

// CreateAdmin 创建新的管理员。
func (s *AdminService) CreateAdmin(username, password, email string, role int) (*model.Admin, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	admin := &model.Admin{Username: username, PasswordHash: hash, Email: email, Role: role}
	if err := s.repo.Create(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

// ChangePassword 更新管理员密码。
func (s *AdminService) ChangePassword(id uint, oldPassword, newPassword string) error {
	admin, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if err := utils.CheckPassword(admin.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("old password incorrect")
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	admin.PasswordHash = hash
	return s.repo.Update(admin)
}

// ListAdmins 返回分页数据。
func (s *AdminService) ListAdmins(page, pageSize int) ([]model.Admin, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}

// GetAdminByID 查询管理员。
func (s *AdminService) GetAdminByID(id uint) (*model.Admin, error) {
	return s.repo.GetByID(id)
}

// UpdateAdmin 修改基本信息。
func (s *AdminService) UpdateAdmin(id uint, email string, role int) (*model.Admin, error) {
	admin, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	admin.Email = email
	admin.Role = role
	if err := s.repo.Update(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

// DeleteAdmin 删除管理员，确保至少保留一个超级管理员。
func (s *AdminService) DeleteAdmin(id uint) error {
	admin, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if mapAdminRole(admin.Role) == "super_admin" {
		count, err := s.repo.CountByRole(admin.Role)
		if err != nil {
			return err
		}
		if count <= 1 {
			return fmt.Errorf("cannot delete the last super admin")
		}
	}
	return s.repo.Delete(id)
}

func mapAdminRole(role int) string {
	if role >= 2 {
		return "super_admin"
	}
	return "admin"
}
