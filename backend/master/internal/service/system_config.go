package service

import (
	"fmt"

	"github.com/iwoov/snell-master/backend/master/internal/repository"
	"github.com/sirupsen/logrus"
)

// SnellConfig Snell Server 配置
type SnellConfig struct {
	Version      string            `json:"version"`
	BaseURL      string            `json:"base_url"`
	DownloadURLs map[string]string `json:"download_urls"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	MasterURL        string
	AgentVersion     string
	AgentDownloadURL string
}

// SystemConfigService 系统配置服务
type SystemConfigService struct {
	repo   repository.SystemConfigRepository
	logger *logrus.Logger
}

// NewSystemConfigService 创建系统配置服务实例
func NewSystemConfigService(repo repository.SystemConfigRepository, logger *logrus.Logger) *SystemConfigService {
	return &SystemConfigService{
		repo:   repo,
		logger: logger,
	}
}

// GetSnellConfig 获取 Snell Server 配置
func (s *SystemConfigService) GetSnellConfig() (*SnellConfig, error) {
	configs, err := s.repo.GetByKeys([]string{
		"snell_version",
		"snell_base_url",
		"snell_mirror_url",
	})
	if err != nil {
		s.logger.Errorf("Failed to get Snell configs: %v", err)
		return nil, err
	}

	version := configs["snell_version"]
	baseURL := configs["snell_base_url"]
	mirrorURL := configs["snell_mirror_url"]

	// 优先使用镜像源
	downloadBase := baseURL
	if mirrorURL != "" {
		downloadBase = mirrorURL
	}

	return &SnellConfig{
		Version: version,
		BaseURL: downloadBase,
		DownloadURLs: map[string]string{
			"amd64":   fmt.Sprintf("%s/snell-server-v%s-linux-amd64.zip", downloadBase, version),
			"i386":    fmt.Sprintf("%s/snell-server-v%s-linux-i386.zip", downloadBase, version),
			"aarch64": fmt.Sprintf("%s/snell-server-v%s-linux-aarch64.zip", downloadBase, version),
			"armv7l":  fmt.Sprintf("%s/snell-server-v%s-linux-armv7l.zip", downloadBase, version),
		},
	}, nil
}

// UpdateSnellVersion 更新 Snell Server 版本
func (s *SystemConfigService) UpdateSnellVersion(version string) error {
	if err := s.repo.Set("snell_version", version); err != nil {
		s.logger.Errorf("Failed to update Snell version: %v", err)
		return err
	}
	s.logger.Infof("Snell version updated to: %s", version)
	return nil
}

// UpdateSnellMirror 更新镜像源
func (s *SystemConfigService) UpdateSnellMirror(mirrorURL string) error {
	if err := s.repo.Set("snell_mirror_url", mirrorURL); err != nil {
		s.logger.Errorf("Failed to update Snell mirror URL: %v", err)
		return err
	}
	s.logger.Infof("Snell mirror URL updated to: %s", mirrorURL)
	return nil
}

// GetSystemConfig 获取系统配置（用于生成部署脚本）
func (s *SystemConfigService) GetSystemConfig() (*SystemInfo, error) {
	configs, err := s.repo.GetByKeys([]string{
		"master_url",
		"agent_version",
		"agent_download_url",
	})
	if err != nil {
		s.logger.Errorf("Failed to get system configs: %v", err)
		return nil, err
	}

	return &SystemInfo{
		MasterURL:        configs["master_url"],
		AgentVersion:     configs["agent_version"],
		AgentDownloadURL: configs["agent_download_url"],
	}, nil
}

// UpdateSystemConfig 更新系统配置
func (s *SystemConfigService) UpdateSystemConfig(key, value string) error {
	if err := s.repo.Set(key, value); err != nil {
		s.logger.Errorf("Failed to update config %s: %v", key, err)
		return err
	}
	s.logger.Infof("System config updated: %s = %s", key, value)
	return nil
}

// GetAllConfigs 获取所有配置
func (s *SystemConfigService) GetAllConfigs() ([]interface{}, error) {
	configs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(configs))
	for i, cfg := range configs {
		result[i] = cfg
	}
	return result, nil
}

// BatchUpdateConfigs 批量更新配置
func (s *SystemConfigService) BatchUpdateConfigs(configs map[string]string) error {
	if err := s.repo.BatchSet(configs); err != nil {
		s.logger.Errorf("Failed to batch update configs: %v", err)
		return err
	}
	s.logger.Info("Batch update configs success")
	return nil
}
