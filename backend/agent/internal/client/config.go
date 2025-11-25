package client

import (
	"encoding/json"
	"fmt"
)

// InstanceConfig 表示 Master 下发的实例配置。
type InstanceConfig struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Port     int    `json:"port"`
	PSK      string `json:"psk"`
	Version  int    `json:"version"`
	OBFS     string `json:"obfs,omitempty"`
}

// ConfigResponse 对应配置拉取接口的响应结构。
type ConfigResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Instances []InstanceConfig `json:"instances"`
	} `json:"data"`
}

// FetchConfig 从 Master 拉取实例配置。
func (c *MasterClient) FetchConfig() ([]InstanceConfig, error) {
	respData, err := c.Get("/api/agent/config")
	if err != nil {
		return nil, err
	}

	var resp ConfigResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal config response: %w", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("fetch config failed: %s", resp.Message)
	}

	return resp.Data.Instances, nil
}
