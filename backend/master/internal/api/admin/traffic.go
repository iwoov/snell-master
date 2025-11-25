package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// TrafficHandler 统计接口。
type TrafficHandler struct {
	svc *service.TrafficService
}

// NewTrafficHandler 构造函数。
func NewTrafficHandler(svc *service.TrafficService) *TrafficHandler {
	return &TrafficHandler{svc: svc}
}

// Summary 返回全局统计。
func (h *TrafficHandler) Summary(c *gin.Context) {
	data, err := h.svc.GetTrafficSummary()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, data)
}

// UserRanking 返回用户排行。
func (h *TrafficHandler) UserRanking(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	data, err := h.svc.GetUserRanking(limit)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, data)
}

// LineTrend 返回趋势。
func (h *TrafficHandler) LineTrend(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	data, err := h.svc.GetTrafficTrend(days)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, data)
}

// NodeTraffic 返回节点统计。
func (h *TrafficHandler) NodeTraffic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid node id")
		return
	}
	data, err := h.svc.GetNodeTraffic(uint(id))
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, data)
}
