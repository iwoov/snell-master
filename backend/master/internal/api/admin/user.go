package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/repository"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// UserHandler 处理管理员侧的用户操作。
type UserHandler struct {
	svc *service.UserService
}

// NewUserHandler 构造函数。
func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// List 用户列表。
func (h *UserHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	filter := repository.UserFilter{Keyword: c.Query("keyword")}
	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filter.Status = &status
		}
	}
	users, total, err := h.svc.GetUserList(page, pageSize, filter)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, gin.H{"items": users, "total": total, "page": page, "page_size": pageSize})
}

// Create 创建用户。
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username     string `json:"username" binding:"required"`
		Password     string `json:"password" binding:"required"`
		Email        string `json:"email"`
		TrafficLimit int64  `json:"traffic_limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := h.svc.CreateUser(req.Username, req.Password, req.Email, req.TrafficLimit)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Created(c, user)
}

// Get 返回用户详情。
func (h *UserHandler) Get(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	user, err := h.svc.GetUserByID(uint(userID))
	if err != nil {
		common.Fail(c, http.StatusNotFound, err.Error())
		return
	}
	common.Success(c, user)
}

// Update 更新用户信息。
func (h *UserHandler) Update(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	user, err := h.svc.UpdateUser(uint(userID), updates)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, user)
}

// Delete 删除用户。
func (h *UserHandler) Delete(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteUser(uint(userID)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"deleted": userID})
}

// ResetTraffic 重置流量。
func (h *UserHandler) ResetTraffic(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.ResetUserTraffic(uint(userID)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"reset": userID})
}

// UpdateStatus 更新用户状态。
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateUserStatus(uint(userID), req.Status); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"status": req.Status})
}

// AssignNodes 分配节点。
func (h *UserHandler) AssignNodes(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	var req struct {
		NodeIDs []uint `json:"node_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.AssignNodes(uint(userID), req.NodeIDs); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"assigned": req.NodeIDs})
}
