package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// InstanceStatusStopped indicates the instance is not running.
	InstanceStatusStopped = iota
	// InstanceStatusRunning indicates the instance process is active.
	InstanceStatusRunning
	// InstanceStatusError indicates the instance failed to start or crashed.
	InstanceStatusError
)

// Instance 表示单个 Snell 实例及其运行状态。
type Instance struct {
	ID       uint
	UserID   uint
	Username string
	Port     int
	PSK      string
	Version  int
	OBFS     string

	PID        int
	PIDFile    string
	ConfigFile string
	LogFile    string
	Status     int

	LastUpdated time.Time
}

// InstanceManager 维护本地实例及其生命周期。
type InstanceManager struct {
	mu             sync.RWMutex
	instances      map[uint]*Instance
	instanceDir    string
	snellBinary    string
	portRangeStart int
	portRangeEnd   int
}

// NewInstanceManager 创建实例管理器并确保必要目录存在。
func NewInstanceManager(instanceDir, snellBinary string, portStart, portEnd int) *InstanceManager {
	_ = os.MkdirAll(instanceDir, 0o755)
	return &InstanceManager{
		instances:      make(map[uint]*Instance),
		instanceDir:    instanceDir,
		snellBinary:    snellBinary,
		portRangeStart: portStart,
		portRangeEnd:   portEnd,
	}
}

// copyInstances 返回当前实例的浅拷贝，供外部遍历使用。
func (m *InstanceManager) copyInstances() map[uint]*Instance {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[uint]*Instance, len(m.instances))
	for id, inst := range m.instances {
		out[id] = inst
	}
	return out
}

// getInstance 在受保护的情况下返回实例。
func (m *InstanceManager) getInstance(id uint) (*Instance, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	inst, ok := m.instances[id]
	return inst, ok
}

// setInstance 更新或新建实例映射。
func (m *InstanceManager) setInstance(inst *Instance) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.instances[inst.ID] = inst
}

// deleteInstance 从映射中移除实例。
func (m *InstanceManager) deleteInstance(id uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.instances, id)
}

// generateFilePaths 根据实例 ID 生成配置、日志、PID 文件路径。
func (m *InstanceManager) generateFilePaths(id uint) (configPath, logPath, pidPath string) {
	configPath = filepath.Join(m.instanceDir, fmt.Sprintf("instance_%d.conf", id))
	logPath = filepath.Join(m.instanceDir, fmt.Sprintf("instance_%d.log", id))
	pidPath = filepath.Join(m.instanceDir, fmt.Sprintf("instance_%d.pid", id))
	return
}

// GetAllInstances 返回实例集合的切片拷贝。
func (m *InstanceManager) GetAllInstances() []*Instance {
	instances := m.copyInstances()
	list := make([]*Instance, 0, len(instances))
	for _, inst := range instances {
		list = append(list, inst)
	}
	return list
}

// GetRunningCount 统计运行中的实例数量。
func (m *InstanceManager) GetRunningCount() int {
	count := 0
	for _, inst := range m.copyInstances() {
		if m.CheckInstanceStatus(inst) == InstanceStatusRunning {
			count++
		}
	}
	return count
}
