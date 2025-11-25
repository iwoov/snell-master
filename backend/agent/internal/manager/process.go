package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iwoov/snell-master/backend/pkg/logger"
	"github.com/iwoov/snell-master/backend/pkg/utils"
)

// StartInstance 启动 Snell 实例（Systemd 模式）。
func (m *InstanceManager) StartInstance(instance *Instance) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}
	if instance.Port <= 0 {
		return fmt.Errorf("instance port must be greater than zero")
	}
	// Systemd 模式下，端口检查可能需要依赖 systemd 启动失败来反馈，或者保留预检查
	if !utils.IsPortAvailable(instance.Port) {
		// 注意：如果是重启，端口可能被旧进程占用，但 Systemd 会处理重启
		// 这里我们只在非运行状态下检查端口
		if !m.isServiceActive(instance.ID) {
			return fmt.Errorf("port %d is not available", instance.Port)
		}
	}

	configPath, logPath := m.generateFilePaths(instance.ID)
	if _, err := m.generateConfig(instance); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}

	// 生成 Systemd Service 文件
	if _, err := m.generateServiceFile(instance); err != nil {
		return fmt.Errorf("generate service file: %w", err)
	}

	// 启动服务
	if err := m.enableAndStartService(instance.ID); err != nil {
		return fmt.Errorf("start service: %w", err)
	}

	instance.ConfigFile = configPath
	instance.LogFile = logPath
	instance.Status = InstanceStatusRunning

	logger.WithModule("manager").Infof("Instance %d started via Systemd (Port=%d)", instance.ID, instance.Port)
	return nil
}

// StopInstance 停止 Snell 实例（Systemd 模式）。
func (m *InstanceManager) StopInstance(instance *Instance) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}

	if err := m.stopAndDisableService(instance.ID); err != nil {
		return fmt.Errorf("stop service: %w", err)
	}

	instance.Status = InstanceStatusStopped
	logger.WithModule("manager").Infof("Instance %d stopped via Systemd", instance.ID)
	return nil
}

// RestartInstance 重启实例。
func (m *InstanceManager) RestartInstance(instance *Instance) error {
	// Ensure config and service files are up to date and service is enabled
	if err := m.StartInstance(instance); err != nil {
		return err
	}

	// Force restart to apply changes
	serviceName := fmt.Sprintf("snell-instance-%d.service", instance.ID)
	if err := m.systemctl("restart", serviceName); err != nil {
		return fmt.Errorf("restart service: %w", err)
	}
	return nil
}

// CheckInstanceStatus 返回实例运行状态。
func (m *InstanceManager) CheckInstanceStatus(instance *Instance) int {
	if instance == nil {
		return InstanceStatusStopped
	}
	if m.isServiceActive(instance.ID) {
		return InstanceStatusRunning
	}
	return InstanceStatusStopped
}

// IsPortUsedByOther 检测端口是否被其他实例占用。
func (m *InstanceManager) IsPortUsedByOther(instanceID uint, port int) bool {
	for _, inst := range m.copyInstances() {
		if inst.ID != instanceID && inst.Port == port {
			return true
		}
	}
	return false
}
