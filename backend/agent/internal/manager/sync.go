package manager

import (
	"os"

	"github.com/iwoov/snell-master/backend/agent/internal/client"
	"github.com/iwoov/snell-master/backend/pkg/logger"
)

// SyncInstances 根据 Master 下发的配置同步本地实例状态。
func (m *InstanceManager) SyncInstances(remoteInstances []client.InstanceConfig) error {
	log := logger.WithModule("manager")
	log.Info("Starting instance sync...")

	remoteIDs := make(map[uint]client.InstanceConfig, len(remoteInstances))
	for _, inst := range remoteInstances {
		remoteIDs[inst.ID] = inst
	}

	for id, localInst := range m.copyInstances() {
		if _, ok := remoteIDs[id]; !ok {
			log.Infof("Deleting instance %d (not present remotely)", id)
			if err := m.StopInstance(localInst); err != nil {
				log.Errorf("Stop instance %d failed: %v", id, err)
			}
			m.deleteInstanceFiles(localInst)
			m.deleteInstance(id)
		}
	}

	for _, remoteInst := range remoteInstances {
		localInst, exists := m.getInstance(remoteInst.ID)
		if !exists {
			log.Infof("Creating new instance %d", remoteInst.ID)
			newInst := &Instance{
				ID:       remoteInst.ID,
				UserID:   remoteInst.UserID,
				Username: remoteInst.Username,
				Port:     remoteInst.Port,
				PSK:      remoteInst.PSK,
				Version:  remoteInst.Version,
				OBFS:     remoteInst.OBFS,
			}
			if err := m.StartInstance(newInst); err != nil {
				log.Errorf("Start instance %d failed: %v", remoteInst.ID, err)
				newInst.Status = InstanceStatusError
			}
			m.setInstance(newInst)
			continue
		}

		if m.isConfigChanged(localInst, remoteInst) {
			log.Infof("Config changed for instance %d, restarting", remoteInst.ID)
			localInst.Port = remoteInst.Port
			localInst.PSK = remoteInst.PSK
			localInst.Version = remoteInst.Version
			localInst.OBFS = remoteInst.OBFS
			if err := m.RestartInstance(localInst); err != nil {
				log.Errorf("Restart instance %d failed: %v", remoteInst.ID, err)
				localInst.Status = InstanceStatusError
			}
		}
	}

	log.Infof("Instance sync completed. Total: %d", len(m.copyInstances()))
	return nil
}

func (m *InstanceManager) isConfigChanged(local *Instance, remote client.InstanceConfig) bool {
	if local == nil {
		return true
	}
	return local.Port != remote.Port ||
		local.PSK != remote.PSK ||
		local.Version != remote.Version ||
		local.OBFS != remote.OBFS
}

func (m *InstanceManager) deleteInstanceFiles(instance *Instance) {
	if instance == nil {
		return
	}
	for _, path := range []string{instance.ConfigFile, instance.PIDFile, instance.LogFile} {
		if path == "" {
			continue
		}
		_ = os.Remove(path)
	}
}
