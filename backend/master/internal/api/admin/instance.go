package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// InstanceHandler 管理 Snell 实例。
type InstanceHandler struct {
	svc *service.InstanceService
}

// NewInstanceHandler 构造函数。
func NewInstanceHandler(svc *service.InstanceService) *InstanceHandler {
	return &InstanceHandler{svc: svc}
}

// List 返回实例列表。
func (h *InstanceHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	filter := repository.InstanceFilter{}
	if userID := c.Query("user_id"); userID != "" {
		if parsed, err := strconv.Atoi(userID); err == nil {
			uid := uint(parsed)
			filter.UserID = &uid
		}
	}
	if nodeID := c.Query("node_id"); nodeID != "" {
		if parsed, err := strconv.Atoi(nodeID); err == nil {
			nid := uint(parsed)
			filter.NodeID = &nid
		}
	}
	filter.Status = c.Query("status")

	instances, err := h.svc.GetInstanceList(filter)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 返回分页格式
	total := len(instances)
	common.Success(c, gin.H{
		"items":     instances,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Create 创建实例。
func (h *InstanceHandler) Create(c *gin.Context) {
	var req struct {
		UserID  uint   `json:"user_id" binding:"required"`
		NodeID  uint   `json:"node_id" binding:"required"`
		Version int    `json:"version"`
		Obfs    string `json:"obfs"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	inst, err := h.svc.CreateInstance(req.UserID, req.NodeID, req.Version, req.Obfs)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Created(c, inst)
}

// Get 返回详情。
func (h *InstanceHandler) Get(c *gin.Context) {
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
	common.Success(c, inst)
}

// Delete 删除实例。
func (h *InstanceHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteInstance(uint(id)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"deleted": id})
}

// UpdateStatus 更新实例状态。
func (h *InstanceHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateInstanceStatus(uint(id), req.Status); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"status": req.Status})
}

// Restart 请求重启实例。
func (h *InstanceHandler) Restart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.RestartInstance(uint(id)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"restarted": id})
}
