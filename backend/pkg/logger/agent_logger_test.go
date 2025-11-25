package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitAgentLogger(t *testing.T) {
	t.Parallel()

	logPath := filepath.Join(t.TempDir(), "agent.log")
	if _, err := InitAgentLogger("info", "json", logPath); err != nil {
		t.Fatalf("InitAgentLogger() error = %v", err)
	}

	WithModule("test").Info("hello")

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	if !strings.Contains(string(data), "\"module\":\"test\"") {
		t.Fatalf("expected module field in log entry: %s", data)
	}
}

func TestInitAgentLoggerInvalidLevel(t *testing.T) {
	t.Parallel()

	if _, err := InitAgentLogger("invalid", "json", ""); err == nil {
		t.Fatal("expected error for invalid log level")
	}
}
