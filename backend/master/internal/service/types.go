package service

import "time"

// TrafficSummary 表示系统流量聚合数据。
type TrafficSummary struct {
	TotalBytes    int64 `json:"total_bytes"`
	TodayBytes    int64 `json:"today_bytes"`
	RecordCount   int64 `json:"record_count"`
	UserCount     int64 `json:"user_count"`
	InstanceCount int64 `json:"instance_count"`
}

// UserTraffic 表示用户排行。
type UserTraffic struct {
	UserID     uint  `json:"user_id"`
	TotalBytes int64 `json:"total_bytes"`
}

// NodeTraffic 统计单个节点。
type NodeTraffic struct {
	NodeID     uint  `json:"node_id"`
	TotalBytes int64 `json:"total_bytes"`
}

// TrafficTrendPoint 供折线展示。
type TrafficTrendPoint struct {
	Date       time.Time `json:"date"`
	TotalBytes int64     `json:"total_bytes"`
}

// DashboardStats 汇总指标。
type DashboardStats struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	TotalNodes     int64 `json:"total_nodes"`
	OnlineNodes    int64 `json:"online_nodes"`
	TotalInstances int64 `json:"total_instances"`
	TotalTraffic   int64 `json:"total_traffic"`
	TodayTraffic   int64 `json:"today_traffic"`
}
