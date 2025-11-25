package client

import (
	"encoding/json"
	"fmt"
)

// SnellConfig 描述 Snell Server 下载配置。
type SnellConfig struct {
	Version      string            `json:"version"`
	BaseURL      string            `json:"base_url"`
	DownloadURLs map[string]string `json:"download_urls"`
}

type snellConfigResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    *SnellConfig `json:"data"`
}

// GetSnellConfig 从 Master 拉取 Snell Server 下载配置。
func (c *MasterClient) GetSnellConfig() (*SnellConfig, error) {
	respData, err := c.Get("/api/agent/snell-config")
	if err != nil {
		return nil, err
	}

	var resp snellConfigResponse
	if err := json.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal snell config: %w", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get snell config failed: %s", resp.Message)
	}
	if resp.Data == nil {
		return nil, fmt.Errorf("snell config payload is empty")
	}
	return resp.Data, nil
}
