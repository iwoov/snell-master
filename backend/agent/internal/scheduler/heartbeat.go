package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/iwoov/snell-master/backend/agent/internal/client"
	"github.com/iwoov/snell-master/backend/agent/internal/manager"
	"github.com/iwoov/snell-master/backend/agent/internal/monitor"
	"github.com/iwoov/snell-master/backend/pkg/logger"
)

const AgentVersion = "1.0.0"

// HeartbeatScheduler 定期向 Master 上报节点状态。
type HeartbeatScheduler struct {
	masterClient  *client.MasterClient
	instanceMgr   *manager.InstanceManager
	systemMonitor *monitor.SystemMonitor

	stopCh   chan struct{}
	stopped  chan struct{}
	interval time.Duration
	once     sync.Once
}

func NewHeartbeatScheduler(masterClient *client.MasterClient, instanceMgr *manager.InstanceManager, systemMonitor *monitor.SystemMonitor) *HeartbeatScheduler {
	return &HeartbeatScheduler{
		masterClient:  masterClient,
		instanceMgr:   instanceMgr,
		systemMonitor: systemMonitor,
	}
}

// Start 启动心跳调度。
func (s *HeartbeatScheduler) Start(intervalSeconds int) error {
	if s.masterClient == nil || s.instanceMgr == nil || s.systemMonitor == nil {
		return fmt.Errorf("heartbeat scheduler dependencies are nil")
	}
	if s.stopCh != nil {
		return fmt.Errorf("heartbeat scheduler already started")
	}
	if intervalSeconds <= 0 {
		intervalSeconds = 30
	}
	s.interval = time.Duration(intervalSeconds) * time.Second
	s.stopCh = make(chan struct{})
	s.stopped = make(chan struct{})

	go s.run()
	go s.sendHeartbeat()

	logger.WithModule("scheduler").Infof("Heartbeat scheduler started (interval: %ds)", intervalSeconds)
	return nil
}

func (s *HeartbeatScheduler) run() {
	ticker := time.NewTicker(s.interval)
	defer func() {
		ticker.Stop()
		close(s.stopped)
	}()

	for {
		select {
		case <-ticker.C:
			s.sendHeartbeat()
		case <-s.stopCh:
			return
		}
	}
}

func (s *HeartbeatScheduler) sendHeartbeat() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.systemMonitor.Update(ctx); err != nil {
		logger.WithModule("scheduler").Warnf("system monitor update failed: %v", err)
	}

	cpuUsage := s.systemMonitor.CPUUsage()
	memUsage := s.systemMonitor.MemoryUsage()
	instanceCount := s.instanceMgr.GetRunningCount()

	if err := s.masterClient.ReportHeartbeat(cpuUsage, memUsage, instanceCount, AgentVersion); err != nil {
		logger.WithModule("scheduler").Errorf("Report heartbeat failed: %v", err)
		return
	}

	logger.WithModule("scheduler").Debugf("Heartbeat sent (CPU: %d%%, MEM: %d%%, instances: %d)", cpuUsage, memUsage, instanceCount)
}

// Stop 停止心跳调度。
func (s *HeartbeatScheduler) Stop() {
	s.once.Do(func() {
		if s.stopCh == nil {
			return
		}
		close(s.stopCh)
		<-s.stopped
		s.stopCh = nil
		logger.WithModule("scheduler").Info("Heartbeat scheduler stopped")
	})
}
