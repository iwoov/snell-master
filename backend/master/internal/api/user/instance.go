package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/api/middleware"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// InstanceHandler 用户查看实例。
type InstanceHandler struct {
	svc *service.InstanceService
}

// NewInstanceHandler 构造函数。
func NewInstanceHandler(svc *service.InstanceService) *InstanceHandler {
	return &InstanceHandler{svc: svc}
}

// ListMyInstances 返回我的实例。
func (h *InstanceHandler) ListMyInstances(c *gin.Context) {
	userID := middleware.GetUserID(c)
	instances, err := h.svc.GetInstancesByUser(userID)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, instances)
}

// GetMyInstance 返回指定实例。
func (h *InstanceHandler) GetMyInstance(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	inst, err := h.svc.GetInstanceByID(uint(id))
	if err != nil {
		common.Fail(c, http.StatusNotFound, err.Error())
		return
	}
	if inst.UserID != userID {
		common.Fail(c, http.StatusForbidden, "not your instance")
		return
	}
	common.Success(c, inst)
}
