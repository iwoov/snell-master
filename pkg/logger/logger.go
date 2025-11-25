package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/iwoov/snell-master/pkg/config"
)

var std = logrus.New()

// Init 根据配置初始化全局日志器。
func Init(cfg config.LogConfig) (*logrus.Logger, error) {
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %s: %w", cfg.Level, err)
	}

	std.SetLevel(level)
	std.SetReportCaller(true)
	std.SetOutput(os.Stdout)

	switch cfg.Format {
	case "json":
		std.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	case "text":
		std.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
	default:
		return nil, fmt.Errorf("invalid log format %s", cfg.Format)
	}

	return std, nil
}

// WithFields 提供结构化字段能力。
func WithFields(fields logrus.Fields) *logrus.Entry {
	return std.WithFields(fields)
}

func Debug(args ...interface{}) { std.Debug(args...) }
func Info(args ...interface{})  { std.Info(args...) }
func Warn(args ...interface{})  { std.Warn(args...) }
func Error(args ...interface{}) { std.Error(args...) }
func Fatal(args ...interface{}) { std.Fatal(args...) }

func Debugf(format string, args ...interface{}) { std.Debugf(format, args...) }
func Infof(format string, args ...interface{})  { std.Infof(format, args...) }
func Warnf(format string, args ...interface{})  { std.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { std.Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { std.Fatalf(format, args...) }
