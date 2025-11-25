package admin

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/iwoov/snell-master/backend/master/internal/api/common"
	"github.com/iwoov/snell-master/backend/master/internal/service"
)

// SubscribeHandler 管理订阅令牌。
type SubscribeHandler struct {
	svc *service.SubscribeService
}

// NewSubscribeHandler 构造函数。
func NewSubscribeHandler(svc *service.SubscribeService) *SubscribeHandler {
	return &SubscribeHandler{svc: svc}
}

// List 返回所有订阅记录。
func (h *SubscribeHandler) List(c *gin.Context) {
	tokens, err := h.svc.ListTokens()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	common.Success(c, tokens)
}

// Create 创建订阅令牌。
func (h *SubscribeHandler) Create(c *gin.Context) {
	var req struct {
		UserID     uint  `json:"user_id" binding:"required"`
		TemplateID *uint `json:"template_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	token, err := h.svc.CreateToken(req.UserID, req.TemplateID, nil)
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Created(c, token)
}

// Delete 删除令牌。
func (h *SubscribeHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteToken(uint(id)); err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"deleted": id})
}

// Regenerate 重新生成令牌。
func (h *SubscribeHandler) Regenerate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, "invalid user id")
		return
	}
	token, err := h.svc.RegenerateToken(uint(id))
	if err != nil {
		common.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	common.Success(c, gin.H{"token": token})
}
