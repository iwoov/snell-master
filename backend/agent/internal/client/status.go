package client

import (
	"encoding/json"
	"fmt"
)

// InstanceStatus 描述单个实例的运行状态。
type InstanceStatus struct {
	InstanceID uint `json:"instance_id"`
	Status     int  `json:"status"`
}

// StatusReportRequest 批量上报实例状态。
type StatusReportRequest struct {
	Statuses []InstanceStatus `json:"statuses"`
}

// StatusReportResponse 表示状态上报的响应。
type StatusReportResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ReportStatus 向 Master 上报实例状态集合。
func (c *MasterClient) ReportStatus(statuses []InstanceStatus) error {
	req := StatusReportRequest{Statuses: statuses}

	data, err := c.Post("/api/agent/status", req)
	if err != nil {
		return err
	}

	var resp StatusReportResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("unmarshal status response: %w", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("status report failed: %s", resp.Message)
	}

	return nil
}
