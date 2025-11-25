package model

import "time"

// TrafficRecord 记录用户流量。
type TrafficRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	InstanceID    uint      `gorm:"index;not null" json:"instance_id"`
	NodeID        uint      `gorm:"index;not null" json:"node_id"`
	UploadBytes   int64     `gorm:"default:0" json:"upload_bytes"`
	DownloadBytes int64     `gorm:"default:0" json:"download_bytes"`
	TotalBytes    int64     `gorm:"column:bytes_total;default:0" json:"total_bytes"`
	RecordDate    time.Time `gorm:"index" json:"record_date"`
	CreatedAt     time.Time `json:"created_at"`

	User     User          `json:"user,omitempty"`
	Node     Node          `json:"node,omitempty"`
	Instance SnellInstance `json:"instance,omitempty"`
}
