package config

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const defaultAgentConfigPath = "backend/agent/configs/agent.yaml"

// AgentConfig 聚合 Agent 运行时所需的所有配置段。
type AgentConfig struct {
	Agent   AgentSettings   `mapstructure:"agent"`
	Monitor MonitorSettings `mapstructure:"monitor"`
}

// AgentSettings 表示 agent 主体的通用配置。
type AgentSettings struct {
	NodeName              string `mapstructure:"node_name"`
	Location              string `mapstructure:"location"`
	CountryCode           string `mapstructure:"country_code"`
	MasterURL             string `mapstructure:"master_url"`
	APIToken              string `mapstructure:"api_token"`
	InstanceDir           string `mapstructure:"instance_dir"`
	PortRangeStart        int    `mapstructure:"port_range_start"`
	PortRangeEnd          int    `mapstructure:"port_range_end"`
	SnellBinary           string `mapstructure:"snell_binary"`
	HeartbeatInterval     int    `mapstructure:"heartbeat_interval"`
	ConfigSyncInterval    int    `mapstructure:"config_sync_interval"`
	TrafficReportInterval int    `mapstructure:"traffic_report_interval"`
	LogLevel              string `mapstructure:"log_level"`
	LogFormat             string `mapstructure:"log_format"`
	LogFile               string `mapstructure:"log_file"`
}

// MonitorSettings 控制监控模块的开关。
type MonitorSettings struct {
	EnableCPU     bool `mapstructure:"enable_cpu"`
	EnableMemory  bool `mapstructure:"enable_memory"`
	EnableTraffic bool `mapstructure:"enable_traffic"`
}

// LoadAgentConfig 读取 YAML 配置文件并应用环境变量覆盖。
func LoadAgentConfig(path string) (*AgentConfig, error) {
	if path == "" {
		path = defaultAgentConfigPath
	}

	if !filepath.IsAbs(path) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("yaml")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	bindAgentEnv(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read agent config: %w", err)
	}

	cfg := &AgentConfig{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal agent config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate 校验配置是否满足最 basic 的约束。
func (c *AgentConfig) Validate() error {
	if c == nil {
		return fmt.Errorf("agent config is nil")
	}

	agent := c.Agent
	trim := func(v string) string { return strings.TrimSpace(v) }

	switch {
	case trim(agent.NodeName) == "":
		return fmt.Errorf("agent.node_name is required")
	case trim(agent.Location) == "":
		return fmt.Errorf("agent.location is required")
	case trim(agent.CountryCode) == "":
		return fmt.Errorf("agent.country_code is required")
	case trim(agent.MasterURL) == "":
		return fmt.Errorf("agent.master_url is required")
	case trim(agent.APIToken) == "":
		return fmt.Errorf("agent.api_token is required")
	case trim(agent.InstanceDir) == "":
		return fmt.Errorf("agent.instance_dir is required")
	case trim(agent.SnellBinary) == "":
		return fmt.Errorf("agent.snell_binary is required")
	case trim(agent.LogLevel) == "":
		return fmt.Errorf("agent.log_level is required")
	case trim(agent.LogFormat) == "":
		return fmt.Errorf("agent.log_format is required")
	}

	if _, err := url.ParseRequestURI(agent.MasterURL); err != nil {
		return fmt.Errorf("agent.master_url is invalid: %w", err)
	}

	if err := validatePortRange(agent.PortRangeStart, agent.PortRangeEnd); err != nil {
		return err
	}

	if agent.HeartbeatInterval <= 0 {
		return fmt.Errorf("agent.heartbeat_interval must be greater than zero")
	}
	if agent.ConfigSyncInterval <= 0 {
		return fmt.Errorf("agent.config_sync_interval must be greater than zero")
	}
	if agent.TrafficReportInterval <= 0 {
		return fmt.Errorf("agent.traffic_report_interval must be greater than zero")
	}

	if err := validateLogFormat(agent.LogFormat); err != nil {
		return err
	}

	if err := validateLogLevel(agent.LogLevel); err != nil {
		return err
	}

	return nil
}

func validatePortRange(start, end int) error {
	switch {
	case start <= 0 || start >= 65535:
		return fmt.Errorf("agent.port_range_start must be between 1 and 65534")
	case end <= 0 || end > 65535:
		return fmt.Errorf("agent.port_range_end must be between 1 and 65535")
	case start >= end:
		return fmt.Errorf("agent.port_range_end must be greater than start")
	}
	return nil
}

func validateLogFormat(format string) error {
	switch strings.ToLower(format) {
	case "json", "text":
		return nil
	default:
		return fmt.Errorf("unsupported agent.log_format %q", format)
	}
}

func validateLogLevel(level string) error {
	valid := map[string]struct{}{
		"debug": {},
		"info":  {},
		"warn":  {},
		"error": {},
	}
	if _, ok := valid[strings.ToLower(level)]; !ok {
		return fmt.Errorf("unsupported agent.log_level %q", level)
	}
	return nil
}

func bindAgentEnv(v *viper.Viper) {
	mappings := map[string]string{
		"agent.node_name":               "AGENT_NODE_NAME",
		"agent.location":                "AGENT_LOCATION",
		"agent.country_code":            "AGENT_COUNTRY_CODE",
		"agent.master_url":              "AGENT_MASTER_URL",
		"agent.api_token":               "AGENT_API_TOKEN",
		"agent.instance_dir":            "AGENT_INSTANCE_DIR",
		"agent.port_range_start":        "AGENT_PORT_RANGE_START",
		"agent.port_range_end":          "AGENT_PORT_RANGE_END",
		"agent.snell_binary":            "AGENT_SNELL_BINARY",
		"agent.heartbeat_interval":      "AGENT_HEARTBEAT_INTERVAL",
		"agent.config_sync_interval":    "AGENT_CONFIG_SYNC_INTERVAL",
		"agent.traffic_report_interval": "AGENT_TRAFFIC_REPORT_INTERVAL",
		"agent.log_level":               "AGENT_LOG_LEVEL",
		"agent.log_format":              "AGENT_LOG_FORMAT",
		"agent.log_file":                "AGENT_LOG_FILE",
		"monitor.enable_cpu":            "AGENT_MONITOR_ENABLE_CPU",
		"monitor.enable_memory":         "AGENT_MONITOR_ENABLE_MEMORY",
		"monitor.enable_traffic":        "AGENT_MONITOR_ENABLE_TRAFFIC",
	}

	for key, env := range mappings {
		_ = v.BindEnv(key, env)
	}
}
