package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iwoov/snell-master/backend/agent/internal/client"
	"github.com/iwoov/snell-master/backend/agent/internal/manager"
	"github.com/iwoov/snell-master/backend/agent/internal/monitor"
	"github.com/iwoov/snell-master/backend/agent/internal/scheduler"
	agentconfig "github.com/iwoov/snell-master/backend/pkg/config"
	"github.com/iwoov/snell-master/backend/pkg/logger"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "backend/agent/configs/agent.yaml", "path to agent config file")
	flag.Parse()

	cfg, err := agentconfig.LoadAgentConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if _, err := logger.InitAgentLogger(cfg.Agent.LogLevel, cfg.Agent.LogFormat, cfg.Agent.LogFile); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}

	log := logger.WithModule("main")
	log.Infof("Starting Snell Agent node=%s location=%s", cfg.Agent.NodeName, cfg.Agent.Location)

	if err := os.MkdirAll(cfg.Agent.InstanceDir, 0o755); err != nil {
		log.Fatalf("create instance dir: %v", err)
	}

	masterClient := client.NewMasterClient(cfg.Agent.MasterURL, cfg.Agent.APIToken)
	instanceMgr := manager.NewInstanceManager(cfg.Agent.InstanceDir, cfg.Agent.SnellBinary, cfg.Agent.PortRangeStart, cfg.Agent.PortRangeEnd)
	systemMonitor := monitor.NewSystemMonitor()
	trafficMonitor := monitor.NewTrafficMonitor(nil)
	snellInstaller := manager.NewSnellInstaller(cfg.Agent.SnellBinary, masterClient)

	if !snellInstaller.IsInstalled() {
		log.Info("Snell binary not found, fetching config from Master...")
		snellCfg, err := masterClient.GetSnellConfig()
		if err != nil {
			log.Fatalf("get snell config: %v", err)
		}
		if err := snellInstaller.Install(context.Background(), snellCfg); err != nil {
			log.Fatalf("install snell: %v", err)
		}
	} else if version, err := snellInstaller.GetVersion(context.Background()); err == nil {
		log.Infof("Snell already installed: %s", version)
	}

	syncScheduler := scheduler.NewSyncScheduler(masterClient, instanceMgr)
	heartbeatScheduler := scheduler.NewHeartbeatScheduler(masterClient, instanceMgr, systemMonitor)
	trafficScheduler := scheduler.NewTrafficScheduler(masterClient, instanceMgr, trafficMonitor)

	if err := syncScheduler.Start(cfg.Agent.ConfigSyncInterval); err != nil {
		log.Fatalf("start sync scheduler: %v", err)
	}
	if err := heartbeatScheduler.Start(cfg.Agent.HeartbeatInterval); err != nil {
		log.Fatalf("start heartbeat scheduler: %v", err)
	}
	if err := trafficScheduler.Start(cfg.Agent.TrafficReportInterval); err != nil {
		log.Fatalf("start traffic scheduler: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	<-ctx.Done()
	stop()

	log.Info("Shutting down Snell Agent...")
	trafficScheduler.Stop()
	heartbeatScheduler.Stop()
	syncScheduler.Stop()

	for _, inst := range instanceMgr.GetAllInstances() {
		if err := instanceMgr.StopInstance(inst); err != nil {
			log.Errorf("stop instance %d: %v", inst.ID, err)
		}
	}

	log.Info("Snell Agent stopped")
}
