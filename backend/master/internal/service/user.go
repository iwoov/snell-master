package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
	jwtutil "github.com/iwoov/snell-master/pkg/jwt"
	"github.com/iwoov/snell-master/pkg/utils"
)

// UserService 处理用户相关业务。
type UserService struct {
	repo          repository.UserRepository
	logger        *logrus.Logger
	jwtSecret     string
	jwtExpireHour int
}

// NewUserService 返回实例。
func NewUserService(repo repository.UserRepository, logger *logrus.Logger, jwtSecret string, jwtExpireHour int) *UserService {
	return &UserService{repo: repo, logger: logger, jwtSecret: jwtSecret, jwtExpireHour: jwtExpireHour}
}

// CreateUser 新建用户并返回结果。
func (s *UserService) CreateUser(username, password, email string, trafficLimit int64) (*model.User, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username:     username,
		PasswordHash: hash,
		Email:        email,
		TrafficLimit: trafficLimit,
		Status:       1,
	}
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserList 返回分页数据。
func (s *UserService) GetUserList(page, pageSize int, filter repository.UserFilter) ([]model.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	return s.repo.List(page, pageSize, filter)
}

// GetUserByID 查询详情。
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	return s.repo.GetByID(id)
}

// UpdateUser 根据字段更新。
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if email, ok := updates["email"].(string); ok {
		user.Email = email
	}
	if limit, ok := getInt64(updates["traffic_limit"]); ok {
		user.TrafficLimit = limit
	}
	if reset, ok := getInt(updates["reset_day"]); ok {
		user.ResetDay = reset
	}
	if status, ok := getInt(updates["status"]); ok {
		user.Status = status
	}
	if expireVal, ok := updates["expire_at"]; ok {
		switch v := expireVal.(type) {
		case string:
			if v != "" {
				if ts, err := time.Parse(time.RFC3339, v); err == nil {
					user.ExpireAt = &ts
				}
			}
		case time.Time:
			user.ExpireAt = &v
		case *time.Time:
			user.ExpireAt = v
		}
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser 删除用户。
func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

// ResetUserTraffic 清零用户流量。
func (s *UserService) ResetUserTraffic(id uint) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	user.TrafficUsedToday = 0
	user.TrafficUsedMonth = 0
	return s.repo.Update(user)
}

// Login 用户登录。
func (s *UserService) Login(username, password string) (string, *model.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, fmt.Errorf("user not found")
		}
		return "", nil, err
	}
	if user.Status == 0 {
		return "", nil, fmt.Errorf("user disabled")
	}
	if err := utils.CheckPassword(user.PasswordHash, password); err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}
	token, err := jwtutil.GenerateToken(user.ID, user.Username, "user", s.jwtSecret, s.jwtExpireHour)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

// ChangePassword 用户修改密码。
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if err := utils.CheckPassword(user.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("old password incorrect")
	}
	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.repo.Update(user)
}

// AssignNodes 重新分配节点。
func (s *UserService) AssignNodes(userID uint, nodeIDs []uint) error {
	return s.repo.AssignNodes(userID, nodeIDs)
}

// UpdateUserStatus 修改状态。
func (s *UserService) UpdateUserStatus(id uint, status int) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	user.Status = status
	return s.repo.Update(user)
}

// GetUserNodes 返回关联节点。
func (s *UserService) GetUserNodes(userID uint) ([]model.Node, error) {
	return s.repo.GetUserNodes(userID)
}

func getInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	case string:
		if v == "" {
			return 0, false
		}
		if parsed, err := strconv.ParseInt(v, 10, 64); err == nil {
			return parsed, true
		}
	}
	return 0, false
}

func getInt(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case string:
		if v == "" {
			return 0, false
		}
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed, true
		}
	}
	return 0, false
}
