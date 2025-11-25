package manager

// GetUsedPorts 返回当前实例占用的端口列表。
func (m *InstanceManager) GetUsedPorts() []int {
	instances := m.copyInstances()
	ports := make([]int, 0, len(instances))
	for _, inst := range instances {
		ports = append(ports, inst.Port)
	}
	return ports
}

// IsPortUsed 检查端口是否被任意实例使用。
func (m *InstanceManager) IsPortUsed(port int) bool {
	for _, inst := range m.copyInstances() {
		if inst.Port == port {
			return true
		}
	}
	return false
}
