package agent

// ConfigResponse 返回节点实例配置。
type ConfigResponse struct {
	Instances []InstanceConfig `json:"instances"`
}

// InstanceConfig 返回节点实例的配置。
type InstanceConfig struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username,omitempty"`
	Port     int    `json:"port"`
	PSK      string `json:"psk"`
	Version  int    `json:"version"`
	Obfs     string `json:"obfs,omitempty"`
}

// HeartbeatRequest 节点心跳上报。
type HeartbeatRequest struct {
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryUsage   float64 `json:"memory_usage"`
	InstanceCount int     `json:"instance_count"`
	Version       string  `json:"version"`
	Status        string  `json:"status"`
}

// TrafficReportRequest 节点流量上报。
type TrafficReportRequest struct {
	Records []TrafficRecord `json:"records"`
}

// TrafficRecord 单条流量记录。
type TrafficRecord struct {
	UserID     uint  `json:"user_id"`
	InstanceID uint  `json:"instance_id"`
	NodeID     uint  `json:"node_id"`
	Upload     int64 `json:"bytes_upload"`
	Download   int64 `json:"bytes_download"`
}

// StatusReportRequest 实例状态上报。
type StatusReportRequest struct {
	Instances []InstanceStatus `json:"instances"`
}

// InstanceStatus 单个实例状态。
type InstanceStatus struct {
	InstanceID uint   `json:"instance_id"`
	Status     string `json:"status"`
}
