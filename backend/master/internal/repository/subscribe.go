package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/iwoov/snell-master/backend/master/internal/model"
)

// SubscribeRepository 管理订阅令牌。
type SubscribeRepository interface {
	Create(token *model.SubscribeToken) error
	GetByID(id uint) (*model.SubscribeToken, error)
	GetByToken(token string) (*model.SubscribeToken, error)
	GetByUser(userID uint) (*model.SubscribeToken, error)
	Update(token *model.SubscribeToken) error
	Delete(id uint) error
	List() ([]model.SubscribeToken, error)
	IncrementAccess(id uint) error
}

type subscribeRepository struct {
	db *gorm.DB
}

// NewSubscribeRepository 返回实现。
func NewSubscribeRepository(db *gorm.DB) SubscribeRepository {
	return &subscribeRepository{db: db}
}

func (r *subscribeRepository) Create(token *model.SubscribeToken) error {
	return r.db.Create(token).Error
}

func (r *subscribeRepository) GetByID(id uint) (*model.SubscribeToken, error) {
	var token model.SubscribeToken
	if err := r.db.Preload("Template").First(&token, id).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *subscribeRepository) GetByToken(token string) (*model.SubscribeToken, error) {
	var sub model.SubscribeToken
	if err := r.db.Where("token = ?", token).Preload("Template").First(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *subscribeRepository) GetByUser(userID uint) (*model.SubscribeToken, error) {
	var token model.SubscribeToken
	if err := r.db.Where("user_id = ?", userID).Preload("Template").First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *subscribeRepository) Update(token *model.SubscribeToken) error {
	return r.db.Save(token).Error
}

func (r *subscribeRepository) Delete(id uint) error {
	return r.db.Delete(&model.SubscribeToken{}, id).Error
}

func (r *subscribeRepository) List() ([]model.SubscribeToken, error) {
	var tokens []model.SubscribeToken
	if err := r.db.Preload("Template").Order("id DESC").Find(&tokens).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *subscribeRepository) IncrementAccess(id uint) error {
	now := time.Now()
	return r.db.Model(&model.SubscribeToken{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"access_count":   gorm.Expr("access_count + 1"),
			"last_access_at": &now,
		}).Error
}
