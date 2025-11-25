package model

import "time"

// SnellInstance 表示运行在节点上的 Snell 服务实例。
type SnellInstance struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`
	NodeID      uint      `gorm:"index;not null" json:"node_id"`
	Port        int       `gorm:"not null" json:"port"`
	PSK         string    `gorm:"size:255;not null" json:"psk"`
	Version     int       `gorm:"default:4" json:"version"`
	Obfs        string    `gorm:"size:64" json:"obfs"`
	ConfigPath  string    `gorm:"size:255" json:"config_path"`
	ServiceName string    `gorm:"size:128" json:"service_name"`
	Status      string    `gorm:"size:32;default:'stopped'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User User `json:"user,omitempty"`
	Node Node `json:"node,omitempty"`
}
