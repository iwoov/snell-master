package model

import "time"

// User 表示普通用户。
type User struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Username         string     `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash     string     `gorm:"size:255;not null" json:"-"`
	Email            string     `gorm:"size:128" json:"email"`
	TrafficLimit     int64      `gorm:"default:107374182400" json:"traffic_limit"` // 100GB
	TrafficUsedToday int64      `gorm:"default:0" json:"traffic_used_today"`
	TrafficUsedMonth int64      `gorm:"default:0" json:"traffic_used_month"`
	TrafficUsedTotal int64      `gorm:"default:0" json:"traffic_used_total"`
	ResetDay         int        `gorm:"default:1" json:"reset_day"`
	Status           int        `gorm:"default:1" json:"status"` // 0: disabled, 1: active
	ExpireAt         *time.Time `json:"expire_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Nodes     []Node          `gorm:"many2many:user_nodes" json:"nodes,omitempty"`
	Instances []SnellInstance `gorm:"foreignKey:UserID" json:"instances,omitempty"`
}
