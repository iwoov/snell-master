package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 聚合所有可配置项，供其他模块依赖。
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

// ServerConfig 定义 HTTP 服务参数。
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 定义 SQLite 相关配置。
type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

// JWTConfig 定义鉴权 token 参数。
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// LogConfig 用于初始化日志系统。
type LogConfig struct {
	Level       string   `mapstructure:"level"`
	Format      string   `mapstructure:"format"`
	IgnorePaths []string `mapstructure:"ignore_paths"`
}

// CORSConfig 控制跨域策略。
type CORSConfig struct {
	AllowOrigins []string `mapstructure:"allow_origins"`
}

// Load 从指定路径加载配置，并允许被环境变量覆盖。
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("config path is required")
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("SNELL")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate 检查必需字段。
func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("server.host is required")
	}
	if c.Server.Port <= 0 {
		return fmt.Errorf("server.port must be greater than zero")
	}
	if c.Database.Path == "" {
		return fmt.Errorf("database.path is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt.secret is required")
	}
	if c.JWT.ExpireHours <= 0 {
		return fmt.Errorf("jwt.expire_hours must be greater than zero")
	}
	if c.Log.Level == "" {
		return fmt.Errorf("log.level is required")
	}
	if c.Log.Format == "" {
		return fmt.Errorf("log.format is required")
	}
	return nil
}

// Address 返回 server host 与 port 组合。
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
