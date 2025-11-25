package model

import "time"

// Node 表示 Snell 节点。
type Node struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Name           string     `gorm:"uniqueIndex;size:64;not null" json:"name"`
	APIToken       string     `gorm:"uniqueIndex;size:128;not null" json:"api_token"`
	Endpoint       string     `gorm:"size:255;not null" json:"endpoint"`
	Location       string     `gorm:"size:100" json:"location"`
	CountryCode    string     `gorm:"size:8" json:"country_code"`
	Status         string     `gorm:"size:32;default:'offline'" json:"status"`
	CPUUsage       float64    `gorm:"default:0" json:"cpu_usage"`
	MemoryUsage    float64    `gorm:"default:0" json:"memory_usage"`
	DiskUsage      float64    `gorm:"default:0" json:"disk_usage"`
	BandwidthUsage float64    `gorm:"default:0" json:"bandwidth_usage"`
	LastSeenAt     *time.Time `json:"last_seen_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Users     []User          `gorm:"many2many:user_nodes" json:"users,omitempty"`
	Instances []SnellInstance `gorm:"foreignKey:NodeID" json:"instances,omitempty"`
}
