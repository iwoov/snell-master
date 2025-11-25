package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	maxPortNumber         = 65535
	defaultCommandTimeout = 30 * time.Second
	publicIPServiceEnv    = "AGENT_PUBLIC_IP_URL"
	defaultPublicIPURL    = "https://api.ipify.org?format=text"
)

// IsPortAvailable 尝试在本地监听端口以判断是否可用。
func IsPortAvailable(port int) bool {
	if port <= 0 || port > maxPortNumber {
		return false
	}
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// FindAvailablePort 在指定范围内查找一个可用端口。
func FindAvailablePort(start, end int) (int, error) {
	if start <= 0 || start >= maxPortNumber {
		return 0, fmt.Errorf("start port must be between 1 and %d", maxPortNumber-1)
	}
	if end <= 0 || end > maxPortNumber {
		return 0, fmt.Errorf("end port must be between 1 and %d", maxPortNumber)
	}
	if start > end {
		return 0, fmt.Errorf("end port must be greater than or equal to start port")
	}

	for port := start; port <= end; port++ {
		if IsPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port in range %d-%d", start, end)
}

// IsProcessRunning 检查指定 PID 是否仍然存活。
func IsProcessRunning(pid int) bool {
	if pid <= 0 {
		return false
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		// Windows 对 Signal(0) 的实现有限，尝试发送 CTRL_BREAK
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			return !errors.Is(err, os.ErrProcessDone)
		}
		return true
	}

	if err := proc.Signal(syscall.Signal(0)); err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return false
		}
		var errno syscall.Errno
		if errors.As(err, &errno) {
			if errno == syscall.ESRCH {
				return false
			}
			if errno == syscall.EPERM {
				return true
			}
		}
		return false
	}
	return true
}

// ReadPIDFile 读取指定路径的 PID 文件。
func ReadPIDFile(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read pid file: %w", err)
	}
	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("invalid pid value %q: %w", pidStr, err)
	}
	return pid, nil
}

// WritePIDFile 将 PID 写入文件，必要时创建目录。
func WritePIDFile(path string, pid int) error {
	if pid <= 0 {
		return fmt.Errorf("pid must be greater than zero")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create pid directory: %w", err)
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0o644)
}

// RemovePIDFile 删除 PID 文件，忽略不存在的情况。
func RemovePIDFile(path string) error {
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("remove pid file: %w", err)
	}
	return nil
}

// ExecuteCommand 带有超时控制地执行命令，并返回标准输出与标准错误。
func ExecuteCommand(command string, args ...string) (string, error) {
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("command is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultCommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("command %s timed out", command)
	}
	if err != nil {
		return "", fmt.Errorf("command %s failed: %w - %s", command, err, strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

// GetHostname 返回当前系统主机名。
func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("get hostname: %w", err)
	}
	return hostname, nil
}

// GetPublicIP 请求外部服务以获取当前公网 IP，服务端点可通过环境变量覆盖。
func GetPublicIP() (string, error) {
	serviceURL := strings.TrimSpace(os.Getenv(publicIPServiceEnv))
	if serviceURL == "" {
		serviceURL = defaultPublicIPURL
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serviceURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request public ip: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("public ip request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 512))
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	ip := strings.TrimSpace(string(data))
	if ip == "" {
		return "", fmt.Errorf("empty response from public ip service")
	}
	return ip, nil
}
