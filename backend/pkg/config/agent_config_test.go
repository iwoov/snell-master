package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadAgentConfig(t *testing.T) {
	t.Parallel()

	cfgPath := writeTempAgentConfig(t, `agent:
  node_name: test-node
  location: Hong Kong
  country_code: HK
  master_url: https://master.example.com
  api_token: test-token
  instance_dir: /var/lib/snell
  port_range_start: 10000
  port_range_end: 20000
  snell_binary: /usr/local/bin/snell-server
  heartbeat_interval: 30
  config_sync_interval: 60
  traffic_report_interval: 300
  log_level: info
  log_format: json
  log_file: /var/log/snell-agent.log
monitor:
  enable_cpu: true
  enable_memory: true
  enable_traffic: false
`)

	cfg, err := LoadAgentConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadAgentConfig() error = %v", err)
	}

	if cfg.Agent.NodeName != "test-node" {
		t.Fatalf("expected node_name=test-node, got %s", cfg.Agent.NodeName)
	}
	if !cfg.Monitor.EnableCPU || !cfg.Monitor.EnableMemory || cfg.Monitor.EnableTraffic {
		t.Fatalf("monitor flags not parsed correctly: %+v", cfg.Monitor)
	}
}

func TestLoadAgentConfigEnvOverride(t *testing.T) {
	cfgPath := writeTempAgentConfig(t, `agent:
  node_name: test-node
  location: Hong Kong
  country_code: HK
  master_url: https://master.example.com
  api_token: test-token
  instance_dir: /var/lib/snell
  port_range_start: 10000
  port_range_end: 20000
  snell_binary: /usr/local/bin/snell-server
  heartbeat_interval: 30
  config_sync_interval: 60
  traffic_report_interval: 300
  log_level: info
  log_format: json
  log_file: /var/log/snell-agent.log
monitor:
  enable_cpu: true
  enable_memory: true
  enable_traffic: false
`)

	t.Setenv("AGENT_MASTER_URL", "https://override.example.com")

	cfg, err := LoadAgentConfig(cfgPath)
	if err != nil {
		t.Fatalf("LoadAgentConfig() error = %v", err)
	}

	if cfg.Agent.MasterURL != "https://override.example.com" {
		t.Fatalf("expected override master url, got %s", cfg.Agent.MasterURL)
	}
}

func TestLoadAgentConfigValidateError(t *testing.T) {
	t.Parallel()

	cfgPath := writeTempAgentConfig(t, `agent:
  node_name: ""
  location: Hong Kong
  country_code: HK
  master_url: https://master.example.com
  api_token: test-token
  instance_dir: /var/lib/snell
  port_range_start: 10000
  port_range_end: 9999
  snell_binary: /usr/local/bin/snell-server
  heartbeat_interval: 30
  config_sync_interval: 60
  traffic_report_interval: 300
  log_level: info
  log_format: json
  log_file: /var/log/snell-agent.log
monitor:
  enable_cpu: true
  enable_memory: true
  enable_traffic: false
`)

	if _, err := LoadAgentConfig(cfgPath); err == nil {
		t.Fatal("expected validation error but got nil")
	} else if !strings.Contains(err.Error(), "agent.node_name") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func writeTempAgentConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "agent.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}
