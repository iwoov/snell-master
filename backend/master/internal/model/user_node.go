package model

import "time"

// UserNode 关联用户和节点。
type UserNode struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"uniqueIndex:idx_user_node;not null" json:"user_id"`
	NodeID    uint       `gorm:"uniqueIndex:idx_user_node;not null" json:"node_id"`
	Connected *time.Time `json:"connected_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// TableName 显式指定表名。
func (UserNode) TableName() string {
	return "user_nodes"
}
