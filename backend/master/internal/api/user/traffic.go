package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/api/middleware"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// TrafficHandler 用户流量查询。
type TrafficHandler struct {
	svc *service.TrafficService
}

// NewTrafficHandler 构造函数。
func NewTrafficHandler(svc *service.TrafficService) *TrafficHandler {
	return &TrafficHandler{svc: svc}
}

// GetMyTraffic 返回指定时间范围的流量。
func (h *TrafficHandler) GetMyTraffic(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var startPtr, endPtr *time.Time
	if startStr := c.Query("start"); startStr != "" {
		if ts, err := time.Parse(time.RFC3339, startStr); err == nil {
			startPtr = &ts
		}
	}
	if endStr := c.Query("end"); endStr != "" {
		if ts, err := time.Parse(time.RFC3339, endStr); err == nil {
			endPtr = &ts
		}
	}
	records, err := h.svc.GetUserTraffic(userID, startPtr, endPtr)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, records)
}
