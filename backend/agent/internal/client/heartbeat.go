package client

import (
	"encoding/json"
	"fmt"
)

// HeartbeatRequest 包含上报的节点心跳信息。
type HeartbeatRequest struct {
	CPUUsage      int    `json:"cpu_usage"`
	MemoryUsage   int    `json:"memory_usage"`
	InstanceCount int    `json:"instance_count"`
	Version       string `json:"version"`
}

// HeartbeatResponse 表示心跳接口的响应。
type HeartbeatResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ReportHeartbeat 将节点心跳上报给 Master。
func (c *MasterClient) ReportHeartbeat(cpuUsage, memUsage, instanceCount int, version string) error {
	req := HeartbeatRequest{
		CPUUsage:      cpuUsage,
		MemoryUsage:   memUsage,
		InstanceCount: instanceCount,
		Version:       version,
	}

	data, err := c.Post("/api/agent/heartbeat", req)
	if err != nil {
		return err
	}

	var resp HeartbeatResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("unmarshal heartbeat response: %w", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("heartbeat failed: %s", resp.Message)
	}

	return nil
}
