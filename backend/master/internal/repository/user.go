package repository

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// UserFilter 用于列表过滤。
type UserFilter struct {
	Status  *int
	Keyword string
}

// UserRepository 定义用户读写接口。
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	List(page, pageSize int, filter UserFilter) ([]model.User, int64, error)
	Update(user *model.User) error
	Delete(id uint) error
	UpdateTraffic(id uint, traffic int64) error
	ResetDailyTraffic() error
	ResetMonthlyTraffic() error
	GetUsersByStatus(status int) ([]model.User, error)
	AssignNodes(userID uint, nodeIDs []uint) error
	GetUserNodes(userID uint) ([]model.Node, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 构造实例。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.Preload("Nodes").Preload("Instances").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(page, pageSize int, filter UserFilter) ([]model.User, int64, error) {
	query := r.db.Model(&model.User{})
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Keyword != "" {
		like := "%" + filter.Keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var users []model.User
	if err := query.Preload("Nodes").Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", id).Delete(&model.UserNode{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.User{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *userRepository) UpdateTraffic(id uint, traffic int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"traffic_used_today": gorm.Expr("traffic_used_today + ?", traffic),
		"traffic_used_month": gorm.Expr("traffic_used_month + ?", traffic),
		"traffic_used_total": gorm.Expr("traffic_used_total + ?", traffic),
	}).Error
}

func (r *userRepository) ResetDailyTraffic() error {
	return r.db.Model(&model.User{}).Update("traffic_used_today", 0).Error
}

func (r *userRepository) ResetMonthlyTraffic() error {
	return r.db.Model(&model.User{}).Updates(map[string]interface{}{
		"traffic_used_month": 0,
		"traffic_used_today": 0,
	}).Error
}

func (r *userRepository) GetUsersByStatus(status int) ([]model.User, error) {
	var users []model.User
	if err := r.db.Where("status = ?", status).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) AssignNodes(userID uint, nodeIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserNode{}).Error; err != nil {
			return err
		}
		if len(nodeIDs) == 0 {
			return nil
		}

		now := time.Now()
		nodes := make([]model.UserNode, 0, len(nodeIDs))
		for _, id := range nodeIDs {
			nodes = append(nodes, model.UserNode{UserID: userID, NodeID: id, Connected: &now})
		}
		return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&nodes).Error
	})
}

func (r *userRepository) GetUserNodes(userID uint) ([]model.Node, error) {
	var nodes []model.Node
	if err := r.db.Model(&model.User{ID: userID}).Association("Nodes").Find(&nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}
