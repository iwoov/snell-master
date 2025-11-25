package manager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/iwoov/snell-master/backend/pkg/logger"
	"github.com/iwoov/snell-master/backend/pkg/utils"
)

// StartInstance 启动 Snell 实例。
func (m *InstanceManager) StartInstance(instance *Instance) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}
	if instance.Port <= 0 {
		return fmt.Errorf("instance port must be greater than zero")
	}
	if m.IsPortUsedByOther(instance.ID, instance.Port) {
		return fmt.Errorf("port %d already used by another instance", instance.Port)
	}
	if !utils.IsPortAvailable(instance.Port) {
		return fmt.Errorf("port %d is not available", instance.Port)
	}

	configPath, logPath, pidPath := m.generateFilePaths(instance.ID)
	if _, err := m.generateConfig(instance); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return fmt.Errorf("create log dir: %w", err)
	}

	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	defer logFile.Close()

	cmd := exec.Command(m.snellBinary, "-c", configPath)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start snell process: %w", err)
	}

	instance.PID = cmd.Process.Pid
	instance.PIDFile = pidPath
	instance.LogFile = logPath
	if err := utils.WritePIDFile(pidPath, instance.PID); err != nil {
		return err
	}

	instance.Status = InstanceStatusRunning
	logger.WithModule("manager").Infof("Instance %d started (PID=%d, Port=%d)", instance.ID, instance.PID, instance.Port)
	return nil
}

// StopInstance 停止 Snell 实例。
func (m *InstanceManager) StopInstance(instance *Instance) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}
	if instance.PID == 0 {
		return fmt.Errorf("instance %d is not running", instance.ID)
	}

	proc, err := os.FindProcess(instance.PID)
	if err != nil {
		return fmt.Errorf("find process: %w", err)
	}

	if runtime.GOOS == "windows" {
		_ = proc.Kill()
	} else {
		if err := proc.Signal(syscall.SIGTERM); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return fmt.Errorf("terminate process: %w", err)
		}
	}

	waitUntil := time.After(2 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	stopped := false
	for !stopped {
		select {
		case <-waitUntil:
			stopped = true
		case <-ticker.C:
			if !utils.IsProcessRunning(instance.PID) {
				stopped = true
			}
		}
	}

	if utils.IsProcessRunning(instance.PID) {
		if err := proc.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			return fmt.Errorf("force kill process: %w", err)
		}
	}

	_ = utils.RemovePIDFile(instance.PIDFile)
	instance.PID = 0
	instance.Status = InstanceStatusStopped
	logger.WithModule("manager").Infof("Instance %d stopped", instance.ID)
	return nil
}

// RestartInstance 重启实例。
func (m *InstanceManager) RestartInstance(instance *Instance) error {
	if err := m.StopInstance(instance); err != nil {
		logger.WithModule("manager").Warnf("Stop before restart failed: %v", err)
	}
	time.Sleep(time.Second)
	return m.StartInstance(instance)
}

// CheckInstanceStatus 返回实例运行状态。
func (m *InstanceManager) CheckInstanceStatus(instance *Instance) int {
	if instance == nil {
		return InstanceStatusStopped
	}
	if instance.PID == 0 {
		return InstanceStatusStopped
	}
	if utils.IsProcessRunning(instance.PID) {
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
