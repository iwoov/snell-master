package client

import (
	"encoding/json"
	"fmt"
)

// InstanceTraffic 包含单个实例的上下行字节统计。
type InstanceTraffic struct {
	InstanceID    uint  `json:"instance_id"`
	BytesUpload   int64 `json:"bytes_upload"`
	BytesDownload int64 `json:"bytes_download"`
}

// TrafficReportRequest 批量上报实例流量。
type TrafficReportRequest struct {
	Traffic []InstanceTraffic `json:"traffic"`
}

// TrafficReportResponse 表示 Master 对流量报告的响应。
type TrafficReportResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ReportTraffic 批量上报实例流量统计。
func (c *MasterClient) ReportTraffic(traffic []InstanceTraffic) error {
	req := TrafficReportRequest{Traffic: traffic}

	data, err := c.Post("/api/agent/traffic", req)
	if err != nil {
		return err
	}

	var resp TrafficReportResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return fmt.Errorf("unmarshal traffic response: %w", err)
	}

	if resp.Code != 0 {
		return fmt.Errorf("traffic report failed: %s", resp.Message)
	}

	return nil
}
