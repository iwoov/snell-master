package model

import "time"

// NodeHeartbeat 存储节点心跳。
type NodeHeartbeat struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NodeID        uint      `gorm:"index;not null" json:"node_id"`
	Status        string    `gorm:"size:32;not null" json:"status"`
	Message       string    `gorm:"size:255" json:"message"`
	CPUUsage      float64   `json:"cpu_usage"`
	MemoryUsage   float64   `json:"memory_usage"`
	InstanceCount int       `json:"instance_count"`
	Version       string    `gorm:"size:32" json:"version"`
	CreatedAt     time.Time `json:"created_at"`

	Node Node `json:"node,omitempty"`
}
