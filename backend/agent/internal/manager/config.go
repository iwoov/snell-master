package manager

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// generateConfig 根据实例信息生成 Snell 配置文件。
func (m *InstanceManager) generateConfig(instance *Instance) (string, error) {
	if instance == nil {
		return "", fmt.Errorf("instance is nil")
	}
	path, _, _ := m.generateFilePaths(instance.ID)
	var builder strings.Builder
	builder.WriteString("[snell-server]\n")
	builder.WriteString(fmt.Sprintf("listen = 0.0.0.0:%d\n", instance.Port))
	builder.WriteString(fmt.Sprintf("psk = %s\n", instance.PSK))
	if strings.TrimSpace(instance.OBFS) != "" {
		builder.WriteString(fmt.Sprintf("obfs = %s\n", instance.OBFS))
	}

	if err := os.WriteFile(path, []byte(builder.String()), 0o644); err != nil {
		return "", fmt.Errorf("write config: %w", err)
	}

	instance.ConfigFile = path
	instance.LastUpdated = time.Now()
	return path, nil
}
