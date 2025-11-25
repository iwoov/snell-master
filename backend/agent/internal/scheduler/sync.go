package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/iwoov/snell-master/backend/agent/internal/client"
	"github.com/iwoov/snell-master/backend/agent/internal/manager"
	"github.com/iwoov/snell-master/backend/pkg/logger"
)

// SyncScheduler 定期与 Master 同步实例配置。
type SyncScheduler struct {
	masterClient *client.MasterClient
	instanceMgr  *manager.InstanceManager

	interval time.Duration
	stopCh   chan struct{}
	stopped  chan struct{}
	once     sync.Once
}

func NewSyncScheduler(masterClient *client.MasterClient, instanceMgr *manager.InstanceManager) *SyncScheduler {
	return &SyncScheduler{masterClient: masterClient, instanceMgr: instanceMgr}
}

func (s *SyncScheduler) Start(intervalSeconds int) error {
	if s.masterClient == nil || s.instanceMgr == nil {
		return fmt.Errorf("sync scheduler dependencies are nil")
	}
	if s.stopCh != nil {
		return fmt.Errorf("sync scheduler already started")
	}
	if intervalSeconds <= 0 {
		intervalSeconds = 60
	}
	s.interval = time.Duration(intervalSeconds) * time.Second
	s.stopCh = make(chan struct{})
	s.stopped = make(chan struct{})

	go s.run()
	go s.syncConfig()

	logger.WithModule("scheduler").Infof("Sync scheduler started (interval: %ds)", intervalSeconds)
	return nil
}

func (s *SyncScheduler) run() {
	ticker := time.NewTicker(s.interval)
	defer func() {
		ticker.Stop()
		close(s.stopped)
	}()

	for {
		select {
		case <-ticker.C:
			s.syncConfig()
		case <-s.stopCh:
			return
		}
	}
}

func (s *SyncScheduler) syncConfig() {
	logger.WithModule("scheduler").Debug("Sync scheduler fetching config")
	instances, err := s.masterClient.FetchConfig()
	if err != nil {
		logger.WithModule("scheduler").Errorf("Fetch config failed: %v", err)
		return
	}
	if err := s.instanceMgr.SyncInstances(instances); err != nil {
		logger.WithModule("scheduler").Errorf("Sync instances failed: %v", err)
	}
}

func (s *SyncScheduler) Stop() {
	s.once.Do(func() {
		if s.stopCh == nil {
			return
		}
		close(s.stopCh)
		<-s.stopped
		s.stopCh = nil
		logger.WithModule("scheduler").Info("Sync scheduler stopped")
	})
}
