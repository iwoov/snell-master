package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// DashboardHandler 返回仪表盘数据。
type DashboardHandler struct {
	svc *service.DashboardService
}

// NewDashboardHandler 构造函数。
func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// Stats 返回统计。
func (h *DashboardHandler) Stats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, stats)
}
