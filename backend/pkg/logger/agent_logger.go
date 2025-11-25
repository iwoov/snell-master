package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	agentLogger   = logrus.New()
	agentLoggerMu sync.Mutex
)

// InitAgentLogger 根据配置初始化日志器，支持控制台与文件双输出。
func InitAgentLogger(level, format, logFile string) (*logrus.Logger, error) {
	agentLoggerMu.Lock()
	defer agentLoggerMu.Unlock()

	lvl, err := logrus.ParseLevel(strings.ToLower(strings.TrimSpace(level)))
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}
	agentLogger.SetLevel(lvl)

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		agentLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	case "text":
		agentLogger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: time.RFC3339})
	default:
		return nil, fmt.Errorf("invalid log format %q", format)
	}

	outputs := []io.Writer{os.Stdout}
	if file := strings.TrimSpace(logFile); file != "" {
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			return nil, fmt.Errorf("create log directory: %w", err)
		}
		f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
		outputs = append(outputs, f)
	}

	agentLogger.SetOutput(io.MultiWriter(outputs...))
	return agentLogger, nil
}

// WithModule 为日志输出追加模块字段，方便追踪来源。
func WithModule(module string) *logrus.Entry {
	return agentLogger.WithField("module", module)
}

// AgentLogger 返回当前的 Agent 日志器实例，主要用于测试场景。
func AgentLogger() *logrus.Logger {
	return agentLogger
}
