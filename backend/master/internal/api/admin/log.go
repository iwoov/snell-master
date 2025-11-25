package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// LogHandler 查询操作日志。
type LogHandler struct {
	svc *service.LogService
}

// NewLogHandler 构造函数。
func NewLogHandler(svc *service.LogService) *LogHandler {
	return &LogHandler{svc: svc}
}

// List 返回日志列表。
func (h *LogHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	filters := map[string]interface{}{}
	if adminID := c.Query("admin_id"); adminID != "" {
		if id, err := strconv.Atoi(adminID); err == nil {
			filters["admin_id"] = id
		}
	}
	if action := c.Query("action"); action != "" {
		filters["action"] = action
	}
	if targetType := c.Query("target_type"); targetType != "" {
		filters["target_type"] = targetType
	}
	logs, total, err := h.svc.GetLogs(page, pageSize, filters)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, gin.H{"items": logs, "total": total, "page": page, "page_size": pageSize})
}
