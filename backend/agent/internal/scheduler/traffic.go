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

// TrafficScheduler 定期收集实例流量并上报。
type TrafficScheduler struct {
	masterClient   *client.MasterClient
	instanceMgr    *manager.InstanceManager
	trafficMonitor *monitor.TrafficMonitor

	interval time.Duration
	stopCh   chan struct{}
	stopped  chan struct{}
	once     sync.Once
}

func NewTrafficScheduler(masterClient *client.MasterClient, instanceMgr *manager.InstanceManager, trafficMonitor *monitor.TrafficMonitor) *TrafficScheduler {
	return &TrafficScheduler{masterClient: masterClient, instanceMgr: instanceMgr, trafficMonitor: trafficMonitor}
}

func (s *TrafficScheduler) Start(intervalSeconds int) error {
	if s.masterClient == nil || s.instanceMgr == nil || s.trafficMonitor == nil {
		return fmt.Errorf("traffic scheduler dependencies are nil")
	}
	if s.stopCh != nil {
		return fmt.Errorf("traffic scheduler already started")
	}
	if intervalSeconds <= 0 {
		intervalSeconds = 300
	}
	s.interval = time.Duration(intervalSeconds) * time.Second
	s.stopCh = make(chan struct{})
	s.stopped = make(chan struct{})

	go s.run()

	logger.WithModule("scheduler").Infof("Traffic scheduler started (interval: %ds)", intervalSeconds)
	return nil
}

func (s *TrafficScheduler) run() {
	reportTicker := time.NewTicker(s.interval)
	syncTicker := time.NewTicker(30 * time.Second)
	defer func() {
		reportTicker.Stop()
		syncTicker.Stop()
		close(s.stopped)
	}()

	// Initial sync
	s.syncRules()

	for {
		select {
		case <-reportTicker.C:
			s.reportTraffic()
		case <-syncTicker.C:
			s.syncRules()
		case <-s.stopCh:
			return
		}
	}
}

func (s *TrafficScheduler) syncRules() {
	instances := s.instanceMgr.GetAllInstances()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.trafficMonitor.SyncRules(ctx, instances); err != nil {
		logger.WithModule("scheduler").Warnf("Sync traffic rules failed: %v", err)
	}
}

func (s *TrafficScheduler) reportTraffic() {
	instances := s.instanceMgr.GetAllInstances()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	traffic, err := s.trafficMonitor.UpdateTraffic(ctx, instances)
	if err != nil {
		logger.WithModule("scheduler").Errorf("Update traffic failed: %v", err)
		return
	}
	if len(traffic) == 0 {
		logger.WithModule("scheduler").Debug("No traffic delta to report")
		return
	}
	if err := s.masterClient.ReportTraffic(traffic); err != nil {
		logger.WithModule("scheduler").Errorf("Report traffic failed: %v", err)
		return
	}
	logger.WithModule("scheduler").Debugf("Traffic reported for %d instances", len(traffic))
}

func (s *TrafficScheduler) Stop() {
	s.once.Do(func() {
		if s.stopCh == nil {
			return
		}
		close(s.stopCh)
		<-s.stopped
		s.stopCh = nil
		logger.WithModule("scheduler").Info("Traffic scheduler stopped")
	})
}
