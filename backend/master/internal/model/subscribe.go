package model

import "time"

// SubscribeToken 用户订阅令牌。
type SubscribeToken struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	UserID       uint       `gorm:"index;not null" json:"user_id"`
	Token        string     `gorm:"uniqueIndex;size:128;not null" json:"token"`
	TemplateID   *uint      `json:"template_id"`
	ExpiresAt    *time.Time `json:"expires_at"`
	LastAccessAt *time.Time `json:"last_access_at"`
	AccessCount  int64      `gorm:"default:0" json:"access_count"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	User     User      `json:"user,omitempty"`
	Template *Template `json:"template,omitempty"`
}
