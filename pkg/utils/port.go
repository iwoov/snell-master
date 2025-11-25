package utils

import (
	"fmt"
	"net"
)

const (
	basePort = 40000
	maxPort  = 60000
)

// AllocatePort 根据用户 ID 分配基础端口。
func AllocatePort(userID uint) int {
	return basePort + int(userID)
}

// IsPortAvailable 尝试在本地监听端口以判断是否可用。
func IsPortAvailable(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// FindAvailablePort 在指定范围内查找可用端口。
func FindAvailablePort(start int) (int, error) {
	if start <= 0 {
		start = basePort
	}
	for port := start; port <= maxPort; port++ {
		if IsPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port in range %d-%d", start, maxPort)
}
