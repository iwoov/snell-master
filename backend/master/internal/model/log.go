package model

import "time"

// OperationLog 记录管理员操作。
type OperationLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	AdminID    uint      `gorm:"index" json:"admin_id"`
	Action     string    `gorm:"size:100;not null" json:"action"`
	TargetType string    `gorm:"size:64" json:"target_type"`
	TargetID   *uint     `json:"target_id"`
	Details    string    `gorm:"type:text" json:"details"`
	IPAddress  string    `gorm:"size:64" json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`

	Admin Admin `json:"admin,omitempty"`
}
