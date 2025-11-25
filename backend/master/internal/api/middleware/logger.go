package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/iwoov/snell-master/pkg/logger"
)

// Logger 记录请求日志并支持忽略指定路径。
func Logger(ignorePaths []string) gin.HandlerFunc {
	skip := make(map[string]struct{})
	for _, path := range ignorePaths {
		if path != "" {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		requestPath := c.FullPath()
		if requestPath == "" {
			requestPath = c.Request.URL.Path
		}
		if _, ok := skip[requestPath]; ok {
			return
		}

		latency := time.Since(start)
		status := c.Writer.Status()
		responseSize := c.Writer.Size()
		rawQuery := c.Request.URL.RawQuery
		pathWithQuery := requestPath
		if rawQuery != "" {
			pathWithQuery = strings.Join([]string{requestPath, rawQuery}, "?")
		}

		fields := logrus.Fields{
			"method":     c.Request.Method,
			"path":       pathWithQuery,
			"status":     status,
			"ip":         c.ClientIP(),
			"latency_ms": latency.Milliseconds(),
			"resp_bytes": responseSize,
			"user_agent": c.Request.UserAgent(),
		}

		if c.Request.ContentLength > 0 {
			fields["req_bytes"] = c.Request.ContentLength
		}
		if c.Errors != nil && len(c.Errors) > 0 {
			fields["error"] = c.Errors.String()
		}

		entry := logger.WithFields(fields)
		switch {
		case status >= 500:
			entry.Error("http request")
		case status >= 400:
			entry.Warn("http request")
		default:
			entry.Info("http request")
		}
	}
}
