package public

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthHandler 提供健康检查。
type HealthHandler struct {
	db        *gorm.DB
	startTime time.Time
}

// NewHealthHandler 构造函数。
func NewHealthHandler(db *gorm.DB, startTime time.Time) *HealthHandler {
	return &HealthHandler{db: db, startTime: startTime}
}

// Ping 简单可用性检查。
func (h *HealthHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// Health 返回系统状态。
func (h *HealthHandler) Health(c *gin.Context) {
	status := "ok"
	dbStatus := "ok"
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil {
			dbStatus = "error"
		} else if err = sqlDB.Ping(); err != nil {
			dbStatus = "error"
		}
	} else {
		dbStatus = "error"
	}
	if dbStatus != "ok" {
		status = "degraded"
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      status,
		"db_status":   dbStatus,
		"started_at":  h.startTime.UTC().Format(time.RFC3339),
		"uptime_secs": int(time.Since(h.startTime).Seconds()),
	})
}
